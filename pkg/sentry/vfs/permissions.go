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
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/syserror"
)

// AccessTypes is a bitmask of Unix file permissions.
type AccessTypes uint16

const (
	MayRead  AccessTypes = 4
	MayWrite             = 2
	MayExec              = 1
)

// GenericCheckPermissions checks that creds has the given access rights on a
// file with the given permissions, UID, and GID, subject to the rules of
// fs/namei.c:generic_permission(). isDir is true if the file is a directory.
func GenericCheckPermissions(creds *auth.Credentials, ats AccessTypes, isDir bool, mode uint16, kuid auth.KUID, kgid auth.KGID) error {
	// Check permission bits.
	perms := mode
	if creds.EffectiveKUID == kuid {
		perms >>= 6
	} else if creds.InGroup(kgid) {
		perms >>= 3
	}
	if uint16(ats)&perms == uint16(ats) {
		return nil
	}

	// Caller capabilities require that the file's KUID and KGID are mapped in
	// the caller's user namespace; compare
	// kernel/capability.c:privileged_wrt_inode_uidgid().
	if !kuid.In(creds.UserNamespace).Ok() || !kgid.In(creds.UserNamespace).Ok() {
		return syserror.EACCES
	}
	// CAP_DAC_READ_SEARCH allows the caller to read and search arbitrary
	// directories, and read arbitrary non-directory files.
	if (isDir && (ats&MayWrite == 0)) || ats == MayRead {
		if creds.HasCapability(linux.CAP_DAC_READ_SEARCH) {
			return nil
		}
	}
	// CAP_DAC_OVERRIDE allows arbitrary access to directories, read/write
	// access to non-directory files, and execute access to non-directory files
	// for which at least one execute bit is set.
	if isDir || (ats&MayExec == 0) || (mode&0111 != 0) {
		if creds.HasCapability(linux.CAP_DAC_OVERRIDE) {
			return nil
		}
	}
	return syserror.EACCES
}

// AccessTypesFromOpenFlags returns the access types required for the given
// OpenOptions.Flags.
func AccessTypesFromOpenFlags(flags uint32) AccessTypes {
	switch flags & linux.O_ACCMODE {
	case linux.O_RDONLY:
		if flags&linux.O_TRUNC != 0 {
			return MayRead | MayWrite
		}
		return MayRead
	case linux.O_WRONLY:
		return MayWrite
	default:
		return MayRead | MayWrite
	}
}
