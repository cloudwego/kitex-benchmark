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
	"io"
	"log"
	"sync"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/streamserver"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex/client/streamxclient"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/streamx"
)

func NewKClient(opt *runner.Options) runner.Client {
	klog.SetLevel(klog.LevelWarn)

	c, err := streamserver.NewClient("test.echo.kitex", streamxclient.WithHostPorts(opt.Address))
	if err != nil {
		log.Fatal(err)
	}
	cli := &kClient{
		client: c,
		streampool: &sync.Pool{
			New: func() interface{} {
				ctx := metainfo.WithValue(context.Background(), "header", "hello")
				stream, err := c.Echo(ctx)
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
	client     streamserver.Client
	streampool *sync.Pool
	reqPool    *sync.Pool
}

func (cli *kClient) Send(method, action, msg string) error {
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	st := cli.streampool.Get()
	if st == nil {
		return errors.New("new stream from streampool failed")
	}
	stream, ok := st.(streamx.BidiStreamingClient[echo.Request, echo.Response])
	if !ok {
		return errors.New("new stream error")
	}
	defer cli.streampool.Put(stream)

	ctx := context.Background()
	req.Action = action
	req.Msg = msg
	err := stream.Send(ctx, req)
	if err != nil {
		return err
	}

	resp, err := stream.Recv(ctx)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	return nil
}

// main is use for routing.
func main() {
	runner.Main("KITEX_TTS_MUX", NewKClient)
}
