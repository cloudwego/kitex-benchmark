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
	"time"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8002
)

var recorder = perf.NewRecorder("KITEX-MUX@Server")

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the EchoImpl interface.
func (s *EchoImpl) Send(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	time.Sleep(time.Duration(req.Time) * time.Millisecond)
	resp := runner.ProcessRequest(recorder, req.Action, "")

	return &echo.Response{
		Action: resp.Action,
		Msg:    resp.Msg,
	}, nil
}

func (s *EchoImpl) StreamTest(stream echo.Echo_StreamTestServer) (err error) {
	return nil
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", port)}
	svr := echosvr.NewServer(
		new(EchoImpl),
		server.WithServiceAddr(address),
		server.WithMuxTransport(),
	)

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
