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
	"sync"

	"github.com/cloudwego/kitex/client"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBKiteXClient(opt *runner.Options) runner.Client {
	cli := &pbKitexClient{}
	cli.client = echosvr.MustNewClient("test.echo.kitex-mux",
		client.WithHostPorts(opt.Address),
		// client.WithMuxConnection(opt.PoolSize))
		// netpoll 设计上，2 个 connection 最优
		client.WithMuxConnection(2))
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &echo.Request{}
		},
	}
	return cli
}

type pbKitexClient struct {
	client  echosvr.Client
	reqPool *sync.Pool
}

func (cli *pbKitexClient) Send(method, action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	req.Msg = msg
	req.Action = action

	reply, err := cli.client.Echo(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
