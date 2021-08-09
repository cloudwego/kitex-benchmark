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

	pb "github.com/cloudwego/kitex-benchmark/protobuf/arpc/pb_gen"
	"github.com/cloudwego/kitex-benchmark/protobuf/arpc/pbcodec"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBArpcClient(opt *runner.Options) runner.Client {
	log.SetLevel(log.LevelNone)

	cli := &pbArpcClient{}
	cli.msg = string(opt.Body)
	cli.msgPool = &sync.Pool{
		New: func() interface{} {
			return &pb.ArpcMsg{}
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
	msg     string
	msgPool *sync.Pool
	clipool *arpc.ClientPool
}

func (cli *pbArpcClient) Echo(msg string) (err error) {
	args := cli.msgPool.Get().(*pb.ArpcMsg)
	reply := cli.msgPool.Get().(*pb.ArpcMsg)
	defer cli.msgPool.Put(args)
	defer cli.msgPool.Put(reply)
	args.Msg = msg

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := cli.clipool.Next()
	err = client.CallWith(ctx, "EchoMsg", args, reply)
	return err
}
