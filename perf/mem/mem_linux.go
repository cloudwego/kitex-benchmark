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
	"regexp"
	"strconv"
	"strings"
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

func parseSmaps(data string) Stats {
	lines := strings.Split(data, "\n")
	reg := regexp.MustCompile("(\\w+): +(\\d+) kB")
	stats := Stats{}
	for _, line := range lines {
		params := reg.FindStringSubmatch(line)
		if len(params) < 3 {
			continue
		}
		field, value := params[1], params[2]
		var err error
		switch strings.ToLower(field) {
		case "rss":
			stats.Rss, err = strconv.ParseInt(value, 10, 64)
		case "pss":
			stats.Pss, err = strconv.ParseInt(value, 10, 64)
		}
		if err != nil {
			panic(fmt.Sprintf("parse mem field: %s failed: %v", field, err))
		}
	}
	return stats
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
