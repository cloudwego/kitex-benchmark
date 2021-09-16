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

	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"

	gogo "github.com/cloudwego/kitex-benchmark/codec/protobuf/gogo_gen"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8003
)

type Echo struct{}

var recorder = perf.NewRecorder("RPCX@Server")

func (s *Echo) Echo(ctx context.Context, args *gogo.Request, reply *gogo.Response) error {
	resp := runner.ProcessRequest(recorder, args.Action, args.Msg)

	reply.Action = resp.Action
	reply.Msg = resp.Msg
	return nil
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	server.UsePool = true
	log.SetDummyLogger()

	s := server.NewServer()
	s.Register(new(Echo), "")
	s.Serve("tcp", fmt.Sprintf(":%d", port))
}
