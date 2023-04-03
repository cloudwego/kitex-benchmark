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

package kclient

import (
	"context"
	"sync"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/generic/data"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewGenericOrdinarySmallClient(opt *runner.Options) runner.Client {
	cli := &genericOrdinaryClient{}
	cli.client = echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.SmallReq
		},
	}
	return cli
}

func NewGenericOrdinaryMediumClient(opt *runner.Options) runner.Client {
	cli := &genericOrdinaryClient{}
	cli.client = echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.MediumReq
		},
	}
	return cli
}

func NewGenericOrdinaryLargeClient(opt *runner.Options) runner.Client {
	cli := &genericOrdinaryClient{}
	cli.client = echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.LargeReq
		},
	}
	return cli
}

type genericOrdinaryClient struct {
	client  echoserver.Client
	reqPool *sync.Pool
}

func (cli *genericOrdinaryClient) Echo(action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.ObjReq)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	reply, err := cli.client.TestObj(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
