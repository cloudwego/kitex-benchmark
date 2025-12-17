/*
 * Copyright 2025 CloudWeGo Authors
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
	"io"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	sechosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/secho"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const port = 8007

var (
	_ echo.SEcho = &EchoImpl{}

	recorder = perf.NewRecorder("KITEX_GRPC_NOREUSE@Server")
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the EchoImpl interface.
func (s *EchoImpl) Echo(ctx context.Context, stream echo.SEcho_EchoServer) error {
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil || len(md["header"]) == 0 || md["header"][0] != "hello" {
		return fmt.Errorf("invalid header: %v", md)
	}

	for {
		req, err := stream.Recv(ctx)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		action, msg := runner.ProcessRequest(recorder, req.Action, req.Msg)

		err = stream.Send(ctx, &echo.Response{
			Action: action,
			Msg:    msg,
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()
	svr := sechosvr.NewServer(
		new(EchoImpl),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4zero, Port: port}),
	)

	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
