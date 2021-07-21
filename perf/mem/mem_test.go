//
// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package mem

import (
	"testing"
)

var testdata = `00400000-7ffc2c73a000 ---p 00000000 00:00 0
[rollup]
Rss:               93112 kB
Pss:               92143 kB
Pss_Anon:          84876 kB
Pss_File:           7267 kB
Pss_Shmem:             0 kB
Shared_Clean:        980 kB
Shared_Dirty:          0 kB
Private_Clean:     61312 kB
Private_Dirty:     30820 kB
Referenced:        39056 kB
Anonymous:         84876 kB
LazyFree:          54056 kB
AnonHugePages:         0 kB
ShmemPmdMapped:        0 kB
FilePmdMapped:         0 kB
Shared_Hugetlb:        0 kB
Private_Hugetlb:       0 kB
Swap:                  0 kB
SwapPss:               0 kB
Locked:                0 kB
`

func TestParseSmaps(t *testing.T) {
	stats := parseSmaps(testdata)
	t.Logf("stats: %v", stats)
	if stats.Rss != 93112 {
		t.Fail()
	}
	if stats.Pss != 92143 {
		t.Fail()
	}
}
