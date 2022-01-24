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
	"fmt"
	"io"
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
		client.WithTransportProtocol(transport.GRPC))
	stream, _ := c.Echo(context.Background())
	cli := &kClient{
		stream: stream,
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &echo.Request{}
			},
		},
		reqQueue:  make(chan *echo.Request, 10240),
		respQueue: make(chan *echo.Response, 10240),
	}
	go func() {
		for req := range cli.reqQueue {
			err := cli.stream.Send(req)
			if err != nil {
				log.Printf("send %s failed: %v\n", req.Action, err)
				cli.respQueue <- nil
			}
		}
	}()
	go func() {
		for {
			resp, err := cli.stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(cli.respQueue)
					close(cli.reqQueue)
					return
				}
				cli.respQueue <- nil
			} else {
				cli.respQueue <- resp
			}
		}
	}()
	return cli
}

type kClient struct {
	stream    sechosvr.SEcho_echoClient
	reqPool   *sync.Pool
	reqQueue  chan *echo.Request
	respQueue chan *echo.Response
}

func (cli *kClient) Echo(action, msg string) error {
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg
	cli.reqQueue <- req

	resp := <-cli.respQueue
	if resp == nil {
		return fmt.Errorf("request fail: %s", req.Action)
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	return nil
}
