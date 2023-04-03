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
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/generic/data"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewGenericMapSmallClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造map 请求和返回类型的泛化调用
	g, err := generic.MapThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericMapClient{}
	cli.client, err = genericclient.NewClient("test.echo.kitex", g,
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	if err != nil {
		panic(err)
	}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.SmallMap
		},
	}
	return cli
}

func NewGenericMapMediumClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造map 请求和返回类型的泛化调用
	g, err := generic.MapThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericMapClient{}
	cli.client, err = genericclient.NewClient("test.echo.kitex", g,
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	if err != nil {
		panic(err)
	}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.MediumMap
		},
	}
	return cli
}

func NewGenericMapLargeClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造map 请求和返回类型的泛化调用
	g, err := generic.MapThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericMapClient{}
	cli.client, err = genericclient.NewClient("test.echo.kitex", g,
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	if err != nil {
		panic(err)
	}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return data.LargeMap
		},
	}
	return cli
}

type genericMapClient struct {
	client  genericclient.Client
	reqPool *sync.Pool
}

func (cli *genericMapClient) Echo(action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(map[string]interface{})
	defer cli.reqPool.Put(req)

	req["action"] = action
	req["msg"] = msg

	reply, err := cli.client.GenericCall(ctx, "TestObj", req)
	if reply != nil {
		repl := reply.(map[string]interface{})
		runner.ProcessResponse(repl["action"].(string), repl["msg"].(string))
	}
	return err
}
