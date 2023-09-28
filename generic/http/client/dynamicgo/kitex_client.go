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
	"bytes"
	"context"
	"fmt"
	"net/http"
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

var (
	p generic.DescriptorProvider
	g generic.Generic
)

func init() {
	var err error
	p, err = generic.NewThriftFileProviderWithDynamicGo("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造http 请求和返回类型的泛化调用
	// enable dynamicgo
	var opts []generic.Option
	opts = append(opts, generic.UseRawBodyForHTTPResp(true))
	g, err = generic.HTTPThriftGeneric(p, opts...)
	if err != nil {
		panic(err)
	}
}

func NewGenericHTTPSmallClient(opt *runner.Options) runner.Client {
	var err error
	cli := &genericHTTPSmallClient{}
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

type genericHTTPSmallClient struct {
	client genericclient.Client
}

func (cli *genericHTTPSmallClient) Echo(action, msg string) error {
	customReq, err := createCustomRequest(action, msg, data.SmallString)
	if err != nil {
		return err
	}

	// send the request
	reply, err := cli.client.GenericCall(context.Background(), "", customReq)
	if reply != nil {
		resp := reply.(*generic.HTTPResponse)
		runner.ProcessResponse(resp.Header.Get("action"), resp.Header.Get("msg"))
	}
	return err
}

func NewGenericHTTPMediumClient(opt *runner.Options) runner.Client {
	var err error
	cli := &genericHTTPMediumClient{}
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

type genericHTTPMediumClient struct {
	client genericclient.Client
}

func (cli *genericHTTPMediumClient) Echo(action, msg string) error {
	customReq, err := createCustomRequest(action, msg, data.MediumString)
	if err != nil {
		return err
	}

	// send the request
	reply, err := cli.client.GenericCall(context.Background(), "", customReq)
	if reply != nil {
		resp := reply.(*generic.HTTPResponse)
		runner.ProcessResponse(resp.Header.Get("action"), resp.Header.Get("msg"))
	}
	return err
}

func NewGenericHTTPLargeClient(opt *runner.Options) runner.Client {
	var err error
	cli := &genericHTTPLargeClient{}
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

type genericHTTPLargeClient struct {
	client genericclient.Client
}

func (cli *genericHTTPLargeClient) Echo(action, msg string) error {
	customReq, err := createCustomRequest(action, msg, data.LargeString)
	if err != nil {
		return err
	}

	// send the request
	reply, err := cli.client.GenericCall(context.Background(), "", customReq)
	if reply != nil {
		resp := reply.(*generic.HTTPResponse)
		runner.ProcessResponse(resp.Header.Get("action"), resp.Header.Get("msg"))
	}
	return err
}

func createCustomRequest(action, msg, data string) (*generic.HTTPRequest, error) {
	url := fmt.Sprintf("http://example.com/test/obj/%s", action)
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("msg", msg)

	customReq, err := generic.FromHTTPRequest(httpRequest)
	if err != nil {
		return nil, err
	}
	return customReq, nil
}
