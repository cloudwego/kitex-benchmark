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

package perf

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"

//	"github.com/felixge/fgprof"
)

func init() {
	mrate, _ := strconv.Atoi(os.Getenv("GOMUTEXRATE"))
	brate, _ := strconv.Atoi(os.Getenv("GOBLOCKRATE"))
	if mrate > 0 {
		runtime.SetMutexProfileFraction(mrate)
	}
	if brate > 0 {
		runtime.SetBlockProfileRate(brate)
	}
}

func ServeMonitor(addr string) error {
//	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	return http.ListenAndServe(addr, nil)
}
