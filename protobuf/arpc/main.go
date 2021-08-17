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
	"github.com/lesismal/arpc/log"

	gogo "github.com/cloudwego/kitex-benchmark/codec/protobuf/gogo_gen"
	"github.com/cloudwego/kitex-benchmark/codec/protobuf/pbcodec"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = ":8004"
)

var recorder = perf.NewRecorder("ARPC@Server")

func Echo(ctx *arpc.Context) {
	args := &gogo.Request{}
	if err := ctx.Bind(args); err != nil {
		ctx.Error(err)
		return
	}

	resp := runner.ProcessRequest(recorder, args.Action, args.Msg)

	reply := &gogo.Response{
		Action: resp.Action,
		Msg:    resp.Msg,
	}
	ctx.Write(reply)
}

func main() {
	log.SetLevel(log.LevelNone)

	codec.DefaultCodec = &pbcodec.ProtoBuffer{}

	svr := arpc.NewServer()
	svr.Handler.EnablePool(true)
	svr.Handler.SetAsyncResponse(true)
	svr.Handler.Handle("Echo", Echo)

	svr.Run(port)
}
