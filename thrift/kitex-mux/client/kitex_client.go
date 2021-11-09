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
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewThriftKiteXClient(opt *runner.Options) runner.Client {
	cli := &thriftKiteXClient{}
	cli.client = echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.Framed),
		client.WithHostPorts(opt.Address),
		client.WithMuxConnection(2))
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &echo.Request{}
		},
	}
	return cli
}

type thriftKiteXClient struct {
	client  echoserver.Client
	reqPool *sync.Pool
}

func (cli *thriftKiteXClient) Echo(action, msg string) error {
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
