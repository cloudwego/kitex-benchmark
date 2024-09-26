/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/cloudwego/kitex-benchmark/runner"
)

// main is use for routing.
func main() {
	if os.Getenv("KITEX_ENABLE_PROFILE") == "1" {
		fmt.Println("[Kitex profile is enabled]")
		// start cpu profile
		cpuProfile, _ := os.Create("output/benchmark-grpc-client-cpu.pprof")
		defer cpuProfile.Close()
		_ = pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()

		// heap profile after finish
		heapProfile, _ := os.Create("output/benchmark-grpc-client-mem.pprof")
		defer func() {
			_ = pprof.WriteHeapProfile(heapProfile)
			heapProfile.Close()
		}()
	}
	runner.Main("KITEX", NewKClient)
}
