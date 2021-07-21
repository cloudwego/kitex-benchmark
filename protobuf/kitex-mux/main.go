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
	"log"
	"net"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/protobuf/kitex/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/protobuf/kitex/kitex_gen/echo/echo"
)

const (
	port = ":8002"
)

var recorder = perf.NewRecorder("KITEX-MUX")

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// EchoStr implements the EchoImpl interface.
func (s *EchoImpl) EchoStr(ctx context.Context, req *echo.StrMsg) (resp *echo.StrMsg, err error) {
	switch req.Msg {
	case "begin":
		recorder.Begin()
	case "end":
		recorder.End()
		recorder.Report()
	}
	return &echo.StrMsg{
		Msg:    req.Msg,
		Finish: req.Finish,
	}, nil
}

func main() {
	address := &net.UnixAddr{Net: "tcp", Name: port}
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
