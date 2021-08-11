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
	"net"
	"sync"
	"time"

	"github.com/lesismal/arpc"
	"github.com/lesismal/arpc/codec"
	"github.com/lesismal/arpc/log"

	gogo "github.com/cloudwego/kitex-benchmark/codec/protobuf/gogo_gen"
	"github.com/cloudwego/kitex-benchmark/codec/protobuf/pbcodec"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBArpcNbioClient(opt *runner.Options) runner.Client {
	log.SetLevel(log.LevelNone)

	cli := &pbArpcClient{}
	cli.msg = string(opt.Body)
	cli.reqPool = &sync.Pool{
		New: func() interface{} {
			return &gogo.Request{}
		},
	}
	cli.respPool = &sync.Pool{
		New: func() interface{} {
			return &gogo.Response{}
		},
	}

	codec.DefaultCodec = &pbcodec.ProtoBuffer{}
	pool, err := arpc.NewClientPool(func() (net.Conn, error) {
		return net.DialTimeout("tcp", opt.Address, time.Second*5)
	}, opt.PoolSize)
	if err != nil {
		panic(err)
	}

	cli.clipool = pool

	return cli
}

type pbArpcClient struct {
	msg      string
	reqPool  *sync.Pool
	respPool *sync.Pool
	clipool  *arpc.ClientPool
}

func (cli *pbArpcClient) Echo(action, msg string) (err error) {
	args := cli.reqPool.Get().(*gogo.Request)
	reply := cli.respPool.Get().(*gogo.Response)
	defer cli.reqPool.Put(args)
	defer cli.respPool.Put(reply)

	args.Action = action
	args.Msg = msg

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client := cli.clipool.Next()
	err = client.CallWith(ctx, "Echo", args, reply)

	runner.ProcessResponse(reply.Action, reply.Msg)
	return err
}
