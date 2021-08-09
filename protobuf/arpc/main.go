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
	"github.com/lesismal/arpc"
	"github.com/lesismal/arpc/codec"

	"github.com/cloudwego/kitex-benchmark/perf"
	pb "github.com/cloudwego/kitex-benchmark/protobuf/arpc/pb_gen"
	"github.com/cloudwego/kitex-benchmark/protobuf/arpc/pbcodec"
)

const (
	port = ":8004"
)

type RpcxEcho struct{}

var recorder = perf.NewRecorder("RPCX")

func EchoMsg(ctx *arpc.Context) {
	args := &pb.ArpcMsg{}
	reply := &pb.ArpcMsg{}
	switch args.Msg {
	case "begin":
		recorder.Begin()
	case "end":
		recorder.End()
		recorder.Report()
	}

	reply.Msg = args.Msg
	ctx.Write(reply)
}

func main() {
	codec.DefaultCodec = &pbcodec.ProtoBuffer{}

	svr := arpc.NewServer()
	svr.Handler.Handle("EchoMsg", EchoMsg)

	svr.Run(port)
}
