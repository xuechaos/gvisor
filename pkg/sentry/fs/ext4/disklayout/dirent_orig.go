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

package disklayout

import "gvisor.dev/gvisor/pkg/sentry/fs"

// DirentOrig represents the original directory entry struct which does not
// contain the file type. This emulates Linux's ext4_dir_entry struct. This is
// used in ext2, ext3 and sometimes in ext4.
//
// Note: This struct can be of variable sizes on disk. The one described below
// is of maximum size and the FileName beyond NameLength bytes might contain
// garbage.
type DirentOrig struct {
	InodeNumber  uint32
	RecordLength uint16
	NameLength   uint16
	FileNameRaw  [MaxFileName]byte
}

// Compiles only if DirentOrig implements Dirent.
var _ Dirent = (*DirentOrig)(nil)

// Inode implements Dirent.Inode.
func (d *DirentOrig) Inode() uint32 { return d.InodeNumber }

// RecordSize implements Dirent.RecordSize.
func (d *DirentOrig) RecordSize() uint16 { return d.RecordLength }

// FileName implements Dirent.FileName.
func (d *DirentOrig) FileName() string {
	return string(d.FileNameRaw[:d.NameLength])
}

// FileType implements Dirent.FileType.
func (d *DirentOrig) FileType() (fs.InodeType, error) { return fs.Anonymous, ErrResolveViaInode }
