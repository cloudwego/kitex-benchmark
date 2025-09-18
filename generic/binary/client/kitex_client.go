/*
 * Copyright 2025 CloudWeGo Authors
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

	"github.com/cloudwego/gopkg/protocol/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
)

var (
	subMsg1 = &echo.SubMessage{
		Id:    &(&struct{ x int64 }{int64(123)}).x,
		Value: &(&struct{ x string }{"hello"}).x,
	}
	subMsg2 = &echo.SubMessage{
		Id:    &(&struct{ x int64 }{int64(321)}).x,
		Value: &(&struct{ x string }{"world"}).x,
	}
	msg1 = &echo.Message{
		Id:          &(&struct{ x int64 }{int64(123)}).x,
		Value:       &(&struct{ x string }{"hello"}).x,
		SubMessages: []*echo.SubMessage{subMsg1, subMsg2},
	}
	msg2 = &echo.Message{
		Id:          &(&struct{ x int64 }{int64(321)}).x,
		Value:       &(&struct{ x string }{"world"}).x,
		SubMessages: []*echo.SubMessage{subMsg2, subMsg1},
	}
)

func NewGenericBinaryClient(opt *runner.Options) runner.Client {
	cli := &genericBinaryClient{}
	var err error
	cli.client, err = genericclient.NewClient("test.echo.kitex", generic.BinaryThriftGenericV2("EchoServer"),
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
			return &echo.ComplexRequest{
				MsgMap: map[string]*echo.SubMessage{
					"v1": subMsg1,
					"v2": subMsg2,
				},
				SubMsgs: []*echo.SubMessage{subMsg1, subMsg2},
				MsgSet:  []*echo.Message{msg1, msg2},
				FlagMsg: msg1,
			}
		},
	}
	return cli
}

type genericBinaryClient struct {
	client  genericclient.Client
	reqPool *sync.Pool
}

func (cli *genericBinaryClient) Send(method, action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.ComplexRequest)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	args := &echo.EchoServerEchoComplexArgs{
		Req: req,
	}
	reply, err := cli.client.GenericCall(ctx, "EchoComplex", thrift.FastMarshal(args))
	if reply != nil {
		result := &echo.EchoServerEchoComplexResult{}
		if err = thrift.FastUnmarshal(reply.([]byte), result); err != nil {
			return err
		}
		runner.ProcessResponse(result.Success.Action, result.Success.Msg)
	}
	return err
}
