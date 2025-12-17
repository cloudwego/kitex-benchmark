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
	"io"
	"log"
	"sync"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote/trans/ttstream"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/streamserver"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewKClient(opt *runner.Options) runner.Client {
	klog.SetLevel(klog.LevelWarn)
	c, err := streamserver.NewClient(
		"test.echo.kitex",
		client.WithHostPorts(opt.Address),
		client.WithTTHeaderStreamingOptions(client.WithTTHeaderStreamingTransportOptions(ttstream.WithClientLongConnPool(ttstream.DefaultLongConnConfig))),
	)
	if err != nil {
		log.Fatal(err)
	}
	cli := &kClient{
		client: c,
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &echo.Request{}
			},
		},
	}
	return cli
}

type kClient struct {
	client  streamserver.Client
	reqPool *sync.Pool
}

func (cli *kClient) Send(method, action, msg string) error {
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	// Create a new stream for each request (no reuse)
	ctx := metainfo.WithValue(context.Background(), "header", "hello")
	stream, err := cli.client.Echo(ctx)
	if err != nil {
		return err
	}
	req.Action = action
	req.Msg = msg
	err = stream.Send(stream.Context(), req)
	if err != nil {
		return err
	}
	err = stream.CloseSend(stream.Context())
	if err != nil {
		return err
	}

	resp, err := stream.Recv(stream.Context())
	if err != nil {
		return err
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	resp, err = stream.Recv(stream.Context())
	if err != io.EOF {
		return err
	}
	return nil
}

// main is use for routing.
func main() {
	runner.Main("KITEX_TTS_LCONN_NOREUSE", NewKClient)
}
