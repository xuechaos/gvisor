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

// Package disklayout provides Linux ext file system's disk level structures
// which can be directly read into from the underlying device. Structs aim to
// emulate structures `exactly` how they are layed out on disk.
//
// This library aims to be compatible with all ext(2/3/4) systems so it
// provides a generic interface for all major structures and various
// implementations (for different versions). The user code is responsible for
// using appropriate implementations based on the underlying device.
//
// Notes:
//   - All fields in these structs are exported because binary.Read would
//     panic otherwise.
//   - All structures on disk are in little-endian order. Only jbd2 (journal)
//     structures are in big-endian order.
//   - All OS dependent fields in these structures will be interpretted using
//     the Linux version of that field.
//   - The suffix `Lo` in field names stands for lower bits of that field.
//   - The suffix `Hi` in field names stands for upper bits of that field. It
//     would need to be left shifted to be used.
//   - The suffix `Raw` has been added to indicate that the field is not split
//     into Lo and Hi fields and also to resolve name collision with the
//     respective interface.
package disklayout
