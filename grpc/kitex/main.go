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
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/pprof"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const port = 8006

var (
	_ echo.Echo = &EchoImpl{}

	recorder = perf.NewRecorder("KITEX@Server")
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the EchoImpl interface.
func (s *EchoImpl) Echo(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	action, msg := runner.ProcessRequest(recorder, req.Action, req.Msg)
	return &echo.Response{
		Action: action,
		Msg:    msg,
	}, nil
}

func main() {
	if os.Getenv("KITEX_ENABLE_PROFILE") == "1" {
		fmt.Println("[Kitex profile is enabled]")
		// start cpu profile
		cpuProfile, _ := os.Create("output/benchmark-grpc-server-cpu.pprof")
		defer cpuProfile.Close()
		_ = pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()

		// heap profile after finish
		heapProfile, _ := os.Create("output/benchmark-grpc-server-mem.pprof")
		defer func() {
			_ = pprof.WriteHeapProfile(heapProfile)
			heapProfile.Close()
		}()
	}
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()
	svr := echosvr.NewServer(
		new(EchoImpl),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4zero, Port: port}),
	)

	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
