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
	"log"
	"net"

	"github.com/cloudwego/gopkg/protocol/thrift"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8004
)

var recorder = perf.NewRecorder("GenericBinary@Server")

// EchoServerImpl implements the last service interface defined in the IDL.
type EchoServerImpl struct{}

func (s *EchoServerImpl) DefaultHandler(ctx context.Context, service, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "Echo":
		args := &echo.EchoServerEchoArgs{}
		if err = thrift.FastUnmarshal(request.([]byte), args); err != nil {
			return nil, err
		}
		req := args.Req
		if req == nil {
			return nil, fmt.Errorf("req is nil")
		}
		action, msg := runner.ProcessRequest(recorder, req.Action, req.Msg)
		resp := &echo.Response{
			Action: action,
			Msg:    msg,
		}
		return thrift.FastMarshal(&echo.EchoServerEchoResult{Success: resp}), nil
	case "EchoComplex":
		args := &echo.EchoServerEchoComplexArgs{}
		if err = thrift.FastUnmarshal(request.([]byte), args); err != nil {
			return nil, err
		}
		req := args.Req
		if req == nil {
			return nil, fmt.Errorf("req is nil")
		}
		action, msg := runner.ProcessRequest(recorder, req.Action, req.Msg)
		return thrift.FastMarshal(&echo.EchoServerEchoComplexResult{Success: &echo.ComplexResponse{
			Action:  action,
			Msg:     msg,
			MsgMap:  req.MsgMap,
			SubMsgs: req.SubMsgs,
			MsgSet:  req.MsgSet,
			FlagMsg: req.FlagMsg,
		}}), nil
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", port)}
	svr := genericserver.NewUnknownServiceOrMethodServer(&genericserver.UnknownServiceOrMethodHandler{
		DefaultHandler: EchoServerImpl{}.DefaultHandler,
	}, server.WithServiceAddr(address), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
