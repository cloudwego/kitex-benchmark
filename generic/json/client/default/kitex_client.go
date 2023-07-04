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
	"encoding/json"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/generic/data"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewGenericJSONSmallClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造json 请求和返回类型的泛化调用
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericJSONSmallClient{}
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
	return cli
}

type genericJSONSmallClient struct {
	client genericclient.Client
}

func (cli *genericJSONSmallClient) Echo(action, msg string) error {
	ctx := context.Background()

	reply, err := cli.client.GenericCall(ctx, "TestObj", data.GetJsonString(action, msg, data.Small))
	if reply != nil {
		repl := reply.(string)
		var rep echo.Request
		err = json.Unmarshal([]byte(repl), &rep)
		if err != nil {
			return err
		}
		runner.ProcessResponse(rep.Action, rep.Msg)
	}
	return err
}

func NewGenericJSONMediumClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造json 请求和返回类型的泛化调用
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericJSONMediumClient{}
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
	return cli
}

type genericJSONMediumClient struct {
	client genericclient.Client
}

func (cli *genericJSONMediumClient) Echo(action, msg string) error {
	ctx := context.Background()

	reply, err := cli.client.GenericCall(ctx, "TestObj", data.GetJsonString(action, msg, data.Medium))
	if reply != nil {
		repl := reply.(string)
		var rep echo.Request
		err = json.Unmarshal([]byte(repl), &rep)
		if err != nil {
			return err
		}
		runner.ProcessResponse(rep.Action, rep.Msg)
	}
	return err
}

func NewGenericJSONLargeClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造json 请求和返回类型的泛化调用
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericJSONLargeClient{}
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
	return cli
}

type genericJSONLargeClient struct {
	client genericclient.Client
}

func (cli *genericJSONLargeClient) Echo(action, msg string) error {
	ctx := context.Background()

	reply, err := cli.client.GenericCall(ctx, "TestObj", data.GetJsonString(action, msg, data.Large))
	if reply != nil {
		repl := reply.(string)
		var rep echo.Request
		err = json.Unmarshal([]byte(repl), &rep)
		if err != nil {
			return err
		}
		runner.ProcessResponse(rep.Action, rep.Msg)
	}
	return err
}
