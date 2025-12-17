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

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/streamserver"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const port = 8004

var (
	_        echo.StreamServer = &StreamServerImpl{}
	recorder                   = perf.NewRecorder("KITEX_TTS_LCONN_NOREUSE@Server")
)

type StreamServerImpl struct{}

func (si *StreamServerImpl) Echo(ctx context.Context, stream echo.StreamServer_EchoServer) error {
	v, _ := metainfo.GetValue(ctx, "header")
	if v != "hello" {
		return fmt.Errorf("invalid header: %v", v)
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

		resp := new(echo.Response)
		resp.Action = action
		resp.Msg = msg
		err = stream.Send(ctx, resp)
		if err != nil {
			return err
		}
	}
}

func main() {
	klog.SetLevel(klog.LevelWarn)
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	svr := streamserver.NewServer(
		new(StreamServerImpl),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4zero, Port: port}),
	)
	if err := svr.Run(); err != nil {
		log.Println(err.Error())
	}
}
