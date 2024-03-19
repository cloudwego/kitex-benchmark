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
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	grpcg "github.com/cloudwego/kitex-benchmark/codec/protobuf/grpc_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBGrpcClient(opt *runner.Options) runner.Client {
	cli := &pbGrpcClient{}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &grpcg.Request{}
		},
	}
	cli.connpool = runner.NewPool(func() interface{} {
		// Set up a connection to the server.
		conn, err := grpc.Dial(opt.Address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return grpcg.NewEchoClient(conn)
	}, opt.PoolSize)
	return cli
}

type pbGrpcClient struct {
	reqPool  *sync.Pool
	connpool *runner.Pool
}

func (cli *pbGrpcClient) Echo(action, msg string) error {
	ctx := context.Background()
	req := cli.reqPool.Get().(*grpcg.Request)
	defer cli.reqPool.Put(req)

	req.Action = action
	req.Msg = msg

	pbcli := cli.connpool.Get().(grpcg.EchoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := pbcli.Echo(ctx, req)

	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}
