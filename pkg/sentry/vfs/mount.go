// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vfs

import (
	"sync/atomic"

	"gvisor.dev/gvisor/pkg/sentry/context"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/syserror"
)

// A Mount is a replacement of a Dentry (Mount.key.point) from one Filesystem
// (Mount.key.parent.fs) with a Dentry (Mount.root) from another Filesystem
// (Mount.fs), which applies to path resolution in the context of a particular
// Mount (Mount.key.parent).
//
// Mounts are reference-counted. Unless otherwise specified, all Mount methods
// require that a reference is held.
//
// Mount and Filesystem are distinct types because it's possible for a single
// Filesystem to be mounted at multiple locations and/or in multiple mount
// namespaces.
//
// Mount is analogous to Linux's struct mount. (gVisor does not distinguish
// between struct mount and struct vfsmount.)
type Mount struct {
	// The lower 63 bits of refs are a reference count. The MSB of refs is set
	// if the Mount has been eagerly unmounted, as by umount(2) without the
	// MNT_DETACH flag. refs is accessed using atomic memory operations.
	refs int64

	// key is protected by VirtualFilesystem.mountMu and
	// VirtualFilesystem.mounts.seq, and may be nil. References are held on
	// key.parent and key.point if they are not nil.
	//
	// Invariant: key.parent != nil iff key.point != nil. key.point belongs to
	// key.parent.fs.
	key mountKey

	// fs, root, and ns are immutable. References are held on fs and root (but
	// not ns).
	//
	// Invariant: root belongs to fs.
	fs   *Filesystem
	root *Dentry
	ns   *MountNamespace
}

// A MountNamespace is a collection of Mounts.
//
// MountNamespaces are reference-counted. Unless otherwise specified, all
// MountNamespace methods require that a reference is held.
//
// MountNamespace is analogous to Linux's struct mnt_namespace.
type MountNamespace struct {
	refs int64 // accessed using atomic memory operations

	// root is the MountNamespace's root mount. root is immutable.
	root *Mount

	// mountpoints contains all Dentries which are mount points in this
	// namespace. mountpoints is protected by VirtualFilesystem.mountMu.
	//
	// mountpoints is used to determine if a Dentry can be moved or removed
	// (which requires that the Dentry is not a mount point in the calling
	// namespace).
	mountpoints map[*Dentry]struct{}
}

// NewMountNamespace returns a new mount namespace with a root filesystem
// configured by the given arguments. A reference is taken on the returned
// MountNamespace.
func (vfs *VirtualFilesystem) NewMountNamespace(ctx context.Context, creds *auth.Credentials, source, fstypename string, opts *NewFilesystemOptions) (*MountNamespace, error) {
	fstype := vfs.getFilesystemType(fstypename)
	if fstype == nil {
		return nil, syserror.ENODEV
	}
	fs, root, err := fstype.NewFilesystem(ctx, creds, source, *opts)
	if err != nil {
		return nil, err
	}
	mntns := &MountNamespace{
		refs:        1,
		mountpoints: make(map[*Dentry]struct{}),
	}
	mntns.root = &Mount{
		fs:   fs,
		root: root,
		ns:   mntns,
		refs: 1,
	}
	return mntns, nil
}

// MountNew creates and mounts a new Filesystem.
func (vfs *VirtualFilesystem) MountNew(ctx context.Context, creds *auth.Credentials, source string, target *PathOperation, fstypename string, opts *NewFilesystemOptions) error {
	fstype := vfs.getFilesystemType(fstypename)
	if fstype == nil {
		return syserror.ENODEV
	}
	fs, root, err := fstype.NewFilesystem(ctx, creds, source, *opts)
	if err != nil {
		return err
	}
	// We can't hold vfs.mountMu while calling FilesystemImpl methods due to
	// lock ordering.
	vd, err := vfs.GetDentryAt(ctx, creds, target, &GetDentryOptions{})
	if err != nil {
		root.decRef(fs)
		fs.decRef()
		return err
	}
	vfs.mountMu.Lock()
	for {
		if vd.dentry.isDisowned() {
			vfs.mountMu.Unlock()
			vd.DecRef()
			root.decRef(fs)
			fs.decRef()
			return syserror.ENOENT
		}
		// vd might have been mounted over between vfs.GetDentryAt() and
		// vfs.mountMu.Lock().
		if !vd.dentry.isMounted() {
			break
		}
		nextmnt := vfs.mounts.Lookup(vd.mount, vd.dentry)
		if nextmnt == nil {
			break
		}
		nextmnt.incRef()
		nextmnt.root.incRef(nextmnt.fs)
		vd.DecRef()
		vd = VirtualDentry{
			mount:  nextmnt,
			dentry: nextmnt.root,
		}
	}
	// TODO: Linux requires that either both the mount point and the mount root
	// are directories, or neither are, and returns ENOTDIR if this is not the
	// case.
	mntns := vd.mount.ns
	mnt := &Mount{
		fs:   fs,
		root: root,
		ns:   mntns,
		refs: 1,
	}
	mnt.storeKey(vd.mount, vd.dentry)
	atomic.AddUint32(&vd.dentry.mounts, 1)
	mntns.mountpoints[vd.dentry] = struct{}{}
	vfsmpmounts, ok := vfs.mountpoints[vd.dentry]
	if !ok {
		vfsmpmounts = make(map[*Mount]struct{})
		vfs.mountpoints[vd.dentry] = vfsmpmounts
	}
	vfsmpmounts[mnt] = struct{}{}
	vfs.mounts.Insert(mnt)
	vfs.mountMu.Unlock()
	return nil
}

// getMountAt returns the last Mount in the stack mounted at (mnt, d). It takes
// a reference on the returned Mount. If (mnt, d) is not a mount point,
// getMountAt returns nil.
//
// getMountAt is analogous to Linux's fs/namei.c:follow_mount().
//
// Preconditions: References are held on mnt and d.
func (vfs *VirtualFilesystem) getMountAt(mnt *Mount, d *Dentry) *Mount {
	// The first mount is special-cased:
	//
	// - The caller is assumed to have checked d.isMounted() already. (This
	// isn't a precondition because it doesn't matter for correctness.)
	//
	// - We return nil, instead of mnt, if there is no mount at (mnt, d).
	//
	// - We don't drop the caller's references on mnt and d.
retryFirst:
	next := vfs.mounts.Lookup(mnt, d)
	if next == nil {
		return nil
	}
	if !next.tryIncMountedRef() {
		// Raced with umount.
		goto retryFirst
	}
	mnt = next
	d = next.root
	// We don't need to take Dentry refs anywhere in this function because
	// Mounts hold references on Mount.root, which is immutable.
	for d.isMounted() {
		next := vfs.mounts.Lookup(mnt, d)
		if next == nil {
			break
		}
		if !next.tryIncMountedRef() {
			// Raced with umount.
			continue
		}
		mnt.decRef()
		mnt = next
		d = next.root
	}
	return mnt
}

// getMountpointAt returns the mount point for the stack of Mounts including
// mnt. It takes a reference on the returned Mount and Dentry. If no such mount
// point exists (i.e. mnt is a root mount), getMountpointAt returns (nil, nil).
//
// Preconditions: References are held on mnt and root. vfsroot is not (mnt,
// mnt.root).
func (vfs *VirtualFilesystem) getMountpointAt(mnt *Mount, vfsroot VirtualDentry) (*Mount, *Dentry) {
	// The first mount is special-cased:
	//
	// - The caller must have already checked mnt against vfsroot.
	//
	// - We return nil, instead of mnt, if there is no mount point for mnt.
	//
	// - We don't drop the caller's reference on mnt.
retryFirst:
	epoch := vfs.mounts.seq.BeginRead()
	parent, point := mnt.loadKey()
	if !vfs.mounts.seq.ReadOk(epoch) {
		goto retryFirst
	}
	if parent == nil {
		return nil, nil
	}
	if !parent.tryIncMountedRef() {
		// Raced with umount.
		goto retryFirst
	}
	if !point.tryIncRef(parent.fs) {
		// Since Mount holds a reference on Mount.key.point, this can only
		// happen due to a racing change to Mount.key.
		parent.decRef()
		goto retryFirst
	}
	mnt = parent
	d := point
	for {
		if mnt == vfsroot.mount && d == vfsroot.dentry {
			break
		}
		if d != mnt.root {
			break
		}
	retryNotFirst:
		epoch := vfs.mounts.seq.BeginRead()
		parent, point := mnt.loadKey()
		if !vfs.mounts.seq.ReadOk(epoch) {
			goto retryNotFirst
		}
		if parent == nil {
			break
		}
		if !parent.tryIncMountedRef() {
			// Raced with umount.
			goto retryNotFirst
		}
		if !point.tryIncRef(parent.fs) {
			// Since Mount holds a reference on Mount.key.point, this can
			// only happen due to a racing change to Mount.key.
			parent.decRef()
			goto retryNotFirst
		}
		if !vfs.mounts.seq.ReadOk(epoch) {
			point.decRef(parent.fs)
			parent.decRef()
			goto retryNotFirst
		}
		d.decRef(mnt.fs)
		mnt.decRef()
		mnt = parent
		d = point
	}
	return mnt, d
}

// tryIncMountedRef increments mnt's reference count and returns true. If mnt's
// reference count is already zero, or has been eagerly unmounted,
// tryIncMountedRef does nothing and returns false.
//
// tryIncMountedRef does not require that a reference is held on mnt.
func (mnt *Mount) tryIncMountedRef() bool {
	for {
		refs := atomic.LoadInt64(&mnt.refs)
		if refs <= 0 { // refs < 0 => MSB set => eagerly unmounted
			return false
		}
		if atomic.CompareAndSwapInt64(&mnt.refs, refs, refs+1) {
			return true
		}
	}
}

func (mnt *Mount) incRef() {
	// In general, negative values for mnt.refs are valid.
	atomic.AddInt64(&mnt.refs, 1)
}

func (mnt *Mount) decRef() {
	refs := atomic.AddInt64(&mnt.refs, -1)
	if refs&int64(^uint64(0)>>1) == 0 { // mask out MSB
		parent, point := mnt.loadKey()
		if point != nil {
			point.decRef(parent.fs)
			parent.decRef()
		}
		mnt.root.decRef(mnt.fs)
		mnt.fs.decRef()
	}
}

// Filesystem returns the mounted Filesystem. It does not take a reference on
// the returned Filesystem.
func (mnt *Mount) Filesystem() *Filesystem {
	return mnt.fs
}

// IncRef increments mntns' reference count.
func (mntns *MountNamespace) IncRef() {
	atomic.AddInt64(&mntns.refs, 1)
}

// DecRef decrements mntns' reference count.
func (mntns *MountNamespace) DecRef() {
	if refs := atomic.AddInt64(&mntns.refs, 0); refs == 0 {
		// TODO: unmount mntns.root
	} else if refs < 0 {
		panic("MountNamespace.DecRef() called without holding a reference")
	}
}

// Root returns mntns' root. A reference is taken on the returned
// VirtualDentry.
func (mntns *MountNamespace) Root() VirtualDentry {
	vd := VirtualDentry{
		mount:  mntns.root,
		dentry: mntns.root.root,
	}
	vd.IncRef()
	return vd
}
