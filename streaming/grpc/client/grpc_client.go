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

	"google.golang.org/grpc"

	grpcg "github.com/cloudwego/kitex-benchmark/codec/protobuf/grpc_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewGrpcClient(opt *runner.Options) runner.Client {
	cli := &grpcClient{
		reqQueue:  make(chan *grpcg.Request, 10240),
		respQueue: make(chan *grpcg.Response, 10240),
	}
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &grpcg.Request{}
		},
	}
	conn, err := grpc.Dial(opt.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := grpcg.NewSEchoClient(conn)
	cli.stream, err = client.Echo(context.Background())
	if err != nil {
		log.Fatalf("did not create stream: %v", err)
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

type grpcClient struct {
	reqPool   *sync.Pool
	stream    grpcg.SEcho_EchoClient
	reqQueue  chan *grpcg.Request
	respQueue chan *grpcg.Response
}

func (cli *grpcClient) Echo(action, msg string) error {
	req := cli.reqPool.Get().(*grpcg.Request)
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
