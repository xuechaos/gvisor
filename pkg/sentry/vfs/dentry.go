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
	"fmt"
	"sync/atomic"
)

// Dentry represents a node in a Filesystem tree which may represent a file.
//
// Dentries are reference-counted. Unless otherwise specified, all Dentry
// methods require that a reference is held.
//
// A Dentry transitions through up to 3 different states through its lifetime:
//
// - Dentries are initially "independent". Independent Dentries have no parent,
// and consequently no name.
//
// - Dentry.InsertChild() causes an independent Dentry to become a "child" of
// another Dentry. A child node has a parent node, and a name in that parent,
// both of which are mutable by DentryMoveChild(). Each child Dentry's name is
// unique within its parent.
//
// - Dentry.RemoveChild() causes a child Dentry to become "disowned". A
// disowned Dentry can still refer to its former parent and its former name in
// said parent, but the disowned Dentry is no longer reachable from its parent,
// and a new Dentry with the same name may become a child of the parent. (This
// is analogous to a struct dentry being "unhashed" in Linux.)
//
// Dentry is loosely analogous to Linux's struct dentry, but:
//
// - VFS does not associate Dentries with inodes. gVisor interacts primarily
// with filesystems that are accessed through filesystem APIs (as opposed to
// raw block devices); many such APIs support only paths and file descriptors,
// and not inodes. Furthermore, when parties outside the scope of VFS can
// rename inodes on such filesystems, VFS generally cannot "follow" the rename,
// both due to synchronization issues and because it may not even be able to
// name the destination path; this implies that it would in fact be *incorrect*
// for Dentries to be associated with inodes on such filesystems. Consequently,
// operations that are inode operations in Linux are FilesystemImpl methods
// and/or FileDescriptionImpl methods in gVisor's VFS. Filesystems that do
// support inodes may store appropriate state in implementations of DentryImpl.
//
// - VFS does not provide synchronization for mutable Dentry fields, other than
// mount-related ones.
//
// - VFS does not require that Dentries are instantiated for all paths accessed
// through VFS, only those that are tracked beyond the scope of a single
// Filesystem operation. This includes file descriptions, mount points, mount
// roots, process working directories, and chroots. This avoids instantiation
// of Dentries for operations on mutable remote filesystems that can't actually
// cache any state in the Dentry.
//
// - For the reasons above, VFS is not directly responsible for managing Dentry
// lifetime. Dentry reference counts only indicate the extent to which VFS
// requires Dentries to exist; Filesystems may elect to cache or discard
// Dentries with zero references.
type Dentry struct {
	// parent is this Dentry's parent in this Filesystem. If this Dentry is
	// independent, parent is nil.
	parent *Dentry

	// name is this Dentry's name in parent.
	name string

	flags uint32

	// mounts is the number of Mounts for which this Dentry is Mount.point.
	// mounts is accessed using atomic memory operations.
	mounts uint32

	// children are child Dentries.
	children map[string]*Dentry

	// impl is the DentryImpl associated with this Dentry. impl is immutable.
	// This should be the last field in Dentry.
	impl DentryImpl
}

const (
	// dflagsDisownedMask is set in Dentry.flags if the Dentry has been
	// disowned.
	dflagsDisownedMask = 1 << iota
)

// Init must be called before first use of d.
func (d *Dentry) Init(impl DentryImpl) {
	d.impl = impl
}

// Impl returns the DentryImpl associated with d.
func (d *Dentry) Impl() DentryImpl {
	return d.impl
}

// DentryImpl contains implementation details for a Dentry. Implementations of
// DentryImpl should contain their associated Dentry by value as their first
// field.
type DentryImpl interface {
	// IncRef increments the Dentry's reference count. A Dentry with a non-zero
	// reference count must remain coherent with the state of the filesystem.
	IncRef(fs *Filesystem)

	// TryIncRef increments the Dentry's reference count and returns true. If
	// the Dentry's reference count is zero, TryIncRef may do nothing and
	// return false. (It is also permitted to succeed if it can restore the
	// guarantee that the Dentry is coherent with the state of the filesystem.)
	//
	// TryIncRef does not require that a reference is held on the Dentry.
	TryIncRef(fs *Filesystem) bool

	// DecRef decrements the Dentry's reference count.
	DecRef(fs *Filesystem)
}

// These functions are exported so that filesystem implementations can use
// them. Users of VFS should not call these functions. Unless otherwise
// specified, all Dentry methods require that there are no concurrent mutators
// of d.

// Name returns d's name in its parent in its owning Filesystem. If d is
// independent, Name returns an empty string.
func (d *Dentry) Name() string {
	return d.name
}

// Parent returns d's parent in its owning Filesystem. It does not take a
// reference on the returned Dentry. If d is independent, Parent returns nil.
func (d *Dentry) Parent() *Dentry {
	return d.parent
}

// ParentOrSelf is equivalent to Parent, but returns d if d is independent.
func (d *Dentry) ParentOrSelf() *Dentry {
	if d.parent == nil {
		return d
	}
	return d.parent
}

// Child returns d's child with the given name in its owning Filesystem. It
// does not take a reference on the returned Dentry. If no such child exists,
// Child returns nil.
func (d *Dentry) Child(name string) *Dentry {
	return d.children[name]
}

// InsertChild makes child a child of d with the given name.
//
// InsertChild is a mutator of d and child.
//
// Preconditions: child must be an independent Dentry. d and child must be from
// the same Filesystem. d must not already have a child with the given name.
func (d *Dentry) InsertChild(child *Dentry, name string) {
	if checkInvariants {
		if _, ok := d.children[name]; ok {
			panic(fmt.Sprintf("parent already contains a child named %q", name))
		}
		if child.parent != nil || child.name != "" {
			panic(fmt.Sprintf("child is not independent: parent = %v, name = %q", child.parent, child.name))
		}
	}
	if d.children == nil {
		d.children = make(map[string]*Dentry)
	}
	d.children[name] = child
	child.parent = d
	child.name = name
}

// TODO: API needs Dentry.RemoveChild{Prepare,Commit,Abort}. Prepare locks
// mountMu for reading and checks that the Dentry to be removed is not a mount
// point in the caller's mount namespace; if it is, Prepare returns EBUSY.
// Commit removes the Dentry, lazily unmounts all mounts at that Dentry, and
// unlocks mountMu; some reference/lock dancing is probably needed to upgrade
// mountMu to writing. Abort just unlocks mountMu. Something similar is needed
// for rename. We might need to introduce something like an inode mutex table
// if filesystems can't safely hold mountMu during remove/rename; Linux uses
// i_mutex, not namespace_sem, for serializing mount activity at a given point,
// but it's not clear if this is necessary for correctness or just better for
// scalability.

func (d *Dentry) isDisowned() bool {
	return atomic.LoadUint32(&d.flags)&dflagsDisownedMask != 0
}

func (d *Dentry) isMounted() bool {
	return atomic.LoadUint32(&d.mounts) != 0
}

func (d *Dentry) incRef(fs *Filesystem) {
	d.impl.IncRef(fs)
}

func (d *Dentry) tryIncRef(fs *Filesystem) bool {
	return d.impl.TryIncRef(fs)
}

func (d *Dentry) decRef(fs *Filesystem) {
	d.impl.DecRef(fs)
}
