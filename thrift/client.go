/*
 * Copyright 2024 CloudWeGo Authors
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

package thrift

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/cloudwego/kitex-tests/pkg/utils"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/runner"
)

type kitexClient struct {
	client echoserver.Client
}

func NewKitexClient(client echoserver.Client) *kitexClient {
	c := new(kitexClient)
	c.client = client
	return c
}

func (cli *kitexClient) Send(method, action, msg string) error {
	switch strings.ToLower(method) {
	case "echo":
		return cli.echo(action, msg)
	case "echocomplex":
		return cli.echoComplex(action, msg)
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
}

var echoReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.Request{}
	},
}

func (cli *kitexClient) echo(action, msg string) error {
	ctx := context.Background()
	req := echoReqPool.Get().(*echo.Request)
	defer echoReqPool.Put(req)

	req.Action = action
	req.Msg = msg

	reply, err := cli.client.Echo(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}

var echoComplexReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.ComplexRequest{}
	},
}

func (cli *kitexClient) echoComplex(action, msg string) error {
	req := createComplexRequest(action, msg)

	reply, err := cli.client.EchoComplex(context.Background(), req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}

func createComplexRequest(action, msg string) *echo.ComplexRequest {
	req := echoComplexReqPool.Get().(*echo.ComplexRequest)
	defer echoComplexReqPool.Put(req)

	id := int64(fastrand.Int31n(100))
	smallSubMsg := &echo.SubMessage{
		Id:    &id,
		Value: stringPtr(utils.RandomString(10)),
	}
	subMsg1K := &echo.SubMessage{
		Id:    &id,
		Value: stringPtr(utils.RandomString(1024)),
	}

	subMsgList2Items := []*echo.SubMessage{smallSubMsg, smallSubMsg}

	message := &echo.Message{
		Id:          &id,
		Value:       stringPtr(utils.RandomString(1024)),
		SubMessages: subMsgList2Items,
	}

	msgMap := make(map[string]*echo.SubMessage)
	for i := 0; i < 5; i++ {
		msgMap[strconv.Itoa(i)] = subMsg1K
	}

	subMsgList100Items := make([]*echo.SubMessage, 100)
	for i := 0; i < len(subMsgList100Items); i++ {
		subMsgList100Items[i] = smallSubMsg
	}

	req.Action = action
	req.Msg = msg
	req.MsgMap = msgMap
	req.SubMsgs = subMsgList100Items
	req.MsgSet = []*echo.Message{message}
	req.FlagMsg = message

	return req
}

func stringPtr(v string) *string { return &v }
