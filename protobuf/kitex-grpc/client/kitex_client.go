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
	"github.com/cloudwego/kitex/transport"
	"sync"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewKClient(opt *runner.Options) runner.Client {
	return &kClient{
		client: echosvr.MustNewClient("test.echo.kitex",
			client.WithHostPorts(opt.Address),
			client.WithLongConnection(
				connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
			client.WithTransportProtocol(transport.GRPC),
		),
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &echo.Request{}
			},
		},
	}
}

type kClient struct {
	client  echosvr.Client
	reqPool *sync.Pool
}

func (cli *kClient) Echo(action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	reply, err := cli.client.Echo(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
