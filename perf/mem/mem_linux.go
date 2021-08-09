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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"unsafe"
)

func getStats(pid int) (Stats, error) {
	data, err := getSmaps(pid)
	if err != nil {
		return Stats{}, err
	}
	return parseSmaps(data), nil
}

func getSmaps(pid int) (string, error) {
	procMemFile := path.Join(fmt.Sprintf("/proc/%d/smaps_rollup", pid))
	file, err := os.Open(procMemFile)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return b2s(bytes), nil
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
