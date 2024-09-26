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
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex-benchmark/thrift"
)

// main is use for routing.
func main() {
	if os.Getenv("KITEX_ENABLE_PROFILE") == "1" {
		fmt.Println("[Kitex profile is enabled]")
		// start cpu profile
		cpuProfile, _ := os.Create("output/benchmark-thrift-client-cpu.pprof")
		defer cpuProfile.Close()
		_ = pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()

		// heap profile after finish
		heapProfile, _ := os.Create("output/benchmark-thrift-client-mem.pprof")
		defer func() {
			_ = pprof.WriteHeapProfile(heapProfile)
			heapProfile.Close()
		}()
	}
	runner.Main("KITEX", NewThriftKitexClient)
}

func NewThriftKitexClient(opt *runner.Options) runner.Client {
	cli := echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.Framed),
		client.WithHostPorts(opt.Address),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	return thrift.NewKitexClient(cli)
}
