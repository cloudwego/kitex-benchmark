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

package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	api "github.com/cloudwego/kitex-benchmark/codec/protobuf/dubbo_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
	"os"
	"sync"
	"time"
)

var echoClientImpl = new(api.EchoClientImpl)

func NewTripleClient(opt *runner.Options) runner.Client {
	cli := &pbTripleClient{}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &api.Request{}
		},
	}

	os.Setenv("DUBBO_GO_CONFIG_PATH", "./protobuf/triple/client/dubbogo.yaml")
	config.SetConsumerService(echoClientImpl)
	if err := config.Load(); err != nil {
		panic(err)
	}

	cli.connpool = runner.NewPool(func() interface{} {
		return echoClientImpl
	}, opt.PoolSize)
	return cli
}

type pbTripleClient struct {
	reqPool  *sync.Pool
	connpool *runner.Pool
}

func (cli *pbTripleClient) Echo(action, msg string) error {
	req := cli.reqPool.Get().(*api.Request)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	tripleCli := cli.connpool.Get().(*api.EchoClientImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err3 := tripleCli.EchoHello(ctx, req)

	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err3
}

func main() {
	runner.Main("tri", NewTripleClient)
}
