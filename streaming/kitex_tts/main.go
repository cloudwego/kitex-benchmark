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
	"io"
	"log"
	"net"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/streamserver"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex/pkg/streamx"
	"github.com/cloudwego/kitex/pkg/streamx/provider/ttstream"
	"github.com/cloudwego/kitex/server"
)

const port = 8002

var (
	_ streamserver.Server = &StreamServerImpl{}

	recorder = perf.NewRecorder("KITEX_TTS@Server")
)

type StreamServerImpl struct{}

func (si *StreamServerImpl) Echo(ctx context.Context, stream streamx.BidiStreamingServer[ttstream.Header, ttstream.Trailer, echo.Request, echo.Response]) error {
	for {
		req, err := stream.Recv(ctx)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		resp := new(echo.Response)
		resp.Msg = req.Msg
		err = stream.Send(ctx, resp)
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

	svr := server.NewServer(server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4zero, Port: port}))
	err := streamserver.RegisterService(svr, new(StreamServerImpl))
	if err != nil {
		panic(err)
	}
	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
