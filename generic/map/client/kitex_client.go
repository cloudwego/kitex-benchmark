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

package main

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

	"github.com/cloudwego/kitex-benchmark/runner"
)

var (
	subMsg1 = map[string]interface{}{
		"id":    int64(123),
		"value": "hello",
	}
	subMsg2 = map[string]interface{}{
		"id":    int64(321),
		"value": "world",
	}
	msg1 = map[string]interface{}{
		"id":          int64(123),
		"value":       "hello",
		"subMessages": []interface{}{subMsg1, subMsg2},
	}
	msg2 = map[string]interface{}{
		"id":          int64(321),
		"value":       "world",
		"subMessages": []interface{}{subMsg2, subMsg1},
	}
)

func NewGenericMapClient(opt *runner.Options) runner.Client {
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
			return map[string]interface{}{
				"msgMap": map[interface{}]interface{}{
					"v1": subMsg1,
					"v2": subMsg2,
				},
				"subMsgs": []interface{}{subMsg1, subMsg2},
				"msgSet":  []interface{}{msg1, msg2},
				"flagMsg": msg1,
			}
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

	reply, err := cli.client.GenericCall(ctx, "EchoComplex", req)
	if reply != nil {
		repl := reply.(map[string]interface{})
		runner.ProcessResponse(repl["action"].(string), repl["msg"].(string))
	}
	return err
}
