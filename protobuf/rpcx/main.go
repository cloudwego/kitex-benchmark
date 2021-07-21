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

	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"

	"github.com/cloudwego/kitex-benchmark/perf"
	pb "github.com/cloudwego/kitex-benchmark/protobuf/rpcx/pb_gen"
)

const (
	port = ":8003"
)

type RpcxEcho struct{}

var recorder = perf.NewRecorder("RPCX")

func (s *RpcxEcho) EchoMsg(ctx context.Context, args *pb.RpcxMsg, reply *pb.RpcxMsg) error {
	switch args.Msg {
	case "begin":
		recorder.Begin()
	case "end":
		recorder.End()
		recorder.Report()
	}

	reply.Msg = args.Msg
	return nil
}

func main() {
	server.UsePool = true
	log.SetDummyLogger()

	s := server.NewServer()
	s.Register(new(RpcxEcho), "")
	s.Serve("tcp", port)
}
