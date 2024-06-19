/*
 * Copyright 2022 CloudWeGo Authors
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

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"

	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8003
)

var recorder = perf.NewRecorder("GenericMap@Server")

type GenericServerImpl struct{}

func (s *GenericServerImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "Echo":
		req := request.(map[string]interface{})
		action, msg := runner.ProcessRequest(recorder, req["action"].(string), req["msg"].(string))
		return map[string]interface{}{
			"action": action,
			"msg":    msg,
		}, nil
	case "EchoComplex":
		req := request.(map[string]interface{})
		action, msg := runner.ProcessRequest(recorder, req["action"].(string), req["msg"].(string))
		return map[string]interface{}{
			"action":  action,
			"msg":     msg,
			"msgMap":  req["msgMap"],
			"subMsgs": req["subMsgs"],
			"msgSet":  req["msgSet"],
			"flagMsg": req["flagMsg"],
		}, nil
	}
	return nil, kerrors.NewBizStatusError(404, "not found")
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	// CurDir: ./scripts
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.MapThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", port)}
	svr := genericserver.NewServer(new(GenericServerImpl), g, server.WithServiceAddr(address), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
