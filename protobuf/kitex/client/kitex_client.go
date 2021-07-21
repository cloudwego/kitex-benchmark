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
	"sync"
	"time"

	"github.com/cloudwego/kitex-benchmark/protobuf/kitex/kitex_gen/echo"
	echo2 "github.com/cloudwego/kitex-benchmark/protobuf/kitex/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
)

func NewPBKiteXClient(opt *runner.Options) runner.Client {
	cli := &pbKitexClient{}
	cli.client = echo2.MustNewClient("test.echo.kitex",
		client.WithHostPorts(opt.Address),
		client.WithLongConnection(connpool.IdleConfig{1000, 1000, time.Minute}))
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &echo.StrMsg{}
		},
	}
	return cli
}

type pbKitexClient struct {
	client  echo2.Client
	reqPool *sync.Pool
}

func (cli *pbKitexClient) Echo(msg string) (err error) {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.StrMsg)
	defer cli.reqPool.Put(req)

	req.Msg = msg
	_, err = cli.client.EchoStr(ctx, req)
	return err
}
