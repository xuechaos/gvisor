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

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/sentry/fs"
)

// Dirent2 represents the ext4 directory entry struct which contains the file
// type. This emulates Linux's ext4_dir_entry_2 struct. The FileName can not be
// more than 255 bytes so we only need 8 bits to store the NameLength. As a
// result, NameLength has been shortened and the rest 8 bits are used to encode
// the file type. Use this only if the SbDirentFileType feature is set.
//
// Note: This struct can be of variable sizes on disk. The one described below
// is of maximum size and the FileName beyond NameLength bytes might contain
// garbage.
type Dirent2 struct {
	InodeNumber  uint32
	RecordLength uint16
	NameLength   uint8
	FileTypeRaw  uint8
	FileNameRaw  [MaxFileName]byte
}

// Compiles only if Dirent2 implements Dirent.
var _ Dirent = (*Dirent2)(nil)

// Inode implements Dirent.Inode.
func (d *Dirent2) Inode() uint32 { return d.InodeNumber }

// RecordSize implements Dirent.RecordSize.
func (d *Dirent2) RecordSize() uint16 { return d.RecordLength }

// FileName implements Dirent.FileName.
func (d *Dirent2) FileName() string {
	return string(d.FileNameRaw[:d.NameLength])
}

// FileType implements Dirent.FileType.
func (d *Dirent2) FileType() (fs.InodeType, error) {
	if inodeType, ok := inodeTypeByFileType[d.FileTypeRaw]; ok {
		return inodeType, nil
	}

	panic(fmt.Sprintf("unknown file type %v", d.FileTypeRaw))
}
