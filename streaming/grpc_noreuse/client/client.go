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

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpcg "github.com/cloudwego/kitex-benchmark/codec/protobuf/grpc_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewGrpcClient(opt *runner.Options) runner.Client {
	conn, err := grpc.Dial(opt.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := grpcg.NewSEchoClient(conn)
	return &grpcClient{
		client: client,
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &grpcg.Request{}
			},
		},
	}
}

type grpcClient struct {
	client  grpcg.SEchoClient
	reqPool *sync.Pool
}

func (cli *grpcClient) Send(method, action, msg string) (err error) {
	req := cli.reqPool.Get().(*grpcg.Request)
	defer cli.reqPool.Put(req)

	// Create a new stream for each request (no reuse)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "header", "hello")
	stream, err := cli.client.Echo(ctx)
	if err != nil {
		return err
	}
	req.Action = action
	req.Msg = msg
	err = stream.Send(req)
	if err != nil {
		return err
	}
	err = stream.CloseSend()
	if err != nil {
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	resp, err = stream.Recv()
	if err != io.EOF {
		return err
	}
	return nil
}

func main() {
	runner.Main("GRPC_NOREUSE", NewGrpcClient)
}
