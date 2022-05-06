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
	"errors"
	"log"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	sechosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/secho"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewKClient(opt *runner.Options) runner.Client {
	klog.SetLevel(klog.LevelWarn)
	c := sechosvr.MustNewClient("test.echo.kitex",
		client.WithHostPorts(opt.Address),
		client.WithTransportProtocol(transport.GRPC),
		client.WithGRPCConnPoolSize(1),
	)
	cli := &kClient{
		client: c,
		streampool: &sync.Pool{
			New: func() interface{} {
				stream, err := c.Echo(context.Background())
				if err != nil {
					log.Printf("client new stream failed: %v", err)
					return nil
				}
				return stream
			},
		},
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &echo.Request{}
			},
		},
	}
	return cli
}

type kClient struct {
	client     sechosvr.Client
	streampool *sync.Pool
	reqPool    *sync.Pool
}

func (cli *kClient) Echo(action, msg string) error {
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	stream, ok := cli.streampool.Get().(sechosvr.SEcho_echoClient)
	if !ok {
		return errors.New("new stream error")
	}
	defer cli.streampool.Put(stream)
	req.Action = action
	req.Msg = msg
	err := stream.Send(req)
	if err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	return nil
}
