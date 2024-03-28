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
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/cloudwego/kitex-benchmark/protobuf/hessian/client/pkg"
	"github.com/cloudwego/kitex-benchmark/runner"
	"os"
	"sync"
)

var (
	echoProvider = &pkg.EchoProvider{}
)

func NewHessianClient(opt *runner.Options) runner.Client {
	cli := &pbTripleClient{}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &pkg.Message{}
		},
	}
	os.Setenv("DUBBO_GO_CONFIG_PATH", "./protobuf/hessian/client/dubbogo.yml")
	hessian.RegisterPOJO(&pkg.Message{})
	config.SetConsumerService(echoProvider)

	if err := config.Load(); err != nil {
		panic(err)
	}

	cli.client = echoProvider
	return cli
}

type pbTripleClient struct {
	reqPool *sync.Pool
	client  *pkg.EchoProvider
}

func (cli *pbTripleClient) Echo(action, msg string) error {
	req := cli.reqPool.Get().(*pkg.Message)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	reply, err := cli.client.Echo(context.Background(), req)

	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
