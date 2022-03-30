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

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"

	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	echosvr "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBKiteXClient(opt *runner.Options) runner.Client {
	cli := &pbKitexClient{}
	cli.client = echosvr.MustNewClient("test.echo.kitex",
		client.WithHostPorts(opt.Address),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &echo.Request{}
		},
	}
	return cli
}

type pbKitexClient struct {
	client  echosvr.Client
	reqPool *sync.Pool
}

func (cli *pbKitexClient) Echo(action, msg string, field, latency, payload int64) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*echo.Request)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Time = latency

	if req.Action == runner.EchoAction {
		if field == 1 {
			req.Field1 = msg
		} else if field == 5 {
			averageLen := (payload) / field
			req.Field1 = msg[0: averageLen]
			req.Field2 = msg[averageLen: 2 * averageLen]
			req.Field3 = msg[averageLen * 2: 3 * averageLen]
			req.Field4 = msg[averageLen * 3: 4 * averageLen]
			req.Field5 = msg[averageLen * 4:]
		} else if field == 10 {
			averageLen := (payload) / field
			req.Field1 = msg[0: averageLen]
			req.Field2 = msg[averageLen: 2 * averageLen]
			req.Field3 = msg[averageLen * 2: 3 * averageLen]
			req.Field4 = msg[averageLen * 3: 4 * averageLen]
			req.Field5 = msg[averageLen * 4: 5 * averageLen]
			req.Field6 = msg[averageLen * 5: 6 * averageLen]
			req.Field7 = msg[averageLen * 6: 7 * averageLen]
			req.Field8 = msg[averageLen * 7: 8 * averageLen]
			req.Field9 = msg[averageLen * 8: 9 * averageLen]
			req.Field10 = msg[averageLen * 9:]
		}
	}

	reply, err := cli.client.Send(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
