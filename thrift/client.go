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

const (
	complexRequestStringFieldCount = 1000
	complexRequestStringFieldSize  = 10 * 1024
)

func (cli *kitexClient) echoComplex(action, msg string) error {
	req := createComplexRequest(action, msg)
	defer echoComplexReqPool.Put(req)

	reply, err := cli.client.EchoComplex(context.Background(), req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}

func createComplexRequest(action, msg string) *echo.ComplexRequest {
	req := echoComplexReqPool.Get().(*echo.ComplexRequest)
	*req = echo.ComplexRequest{}

	id := int64(fastrand.Int31n(100))
	largeValue := utils.RandomString(complexRequestStringFieldSize)

	msgMapSize := complexRequestStringFieldCount / 2
	subMsgListSize := complexRequestStringFieldCount - msgMapSize

	msgMap := make(map[string]*echo.SubMessage, msgMapSize)
	for i := 0; i < msgMapSize; i++ {
		subID := id + int64(i)
		msgMap[strconv.Itoa(i)] = &echo.SubMessage{
			Id:    &subID,
			Value: ptr(largeValue),
		}
	}

	message := &echo.Message{
		Id:          &id,
		Value:       ptr(largeValue),
		SubMessages: nil,
	}

	subMsgs := make([]*echo.SubMessage, subMsgListSize)
	for i := 0; i < len(subMsgs); i++ {
		subID := id + int64(msgMapSize+i)
		subMsgs[i] = &echo.SubMessage{
			Id:    &subID,
			Value: ptr(largeValue),
		}
	}

	req.Action = action
	req.Msg = msg
	req.MsgMap = msgMap
	req.SubMsgs = subMsgs
	req.MsgSet = []*echo.Message{message}
	req.FlagMsg = message

	return req
}

func ptr[T any](v T) *T { return &v }
