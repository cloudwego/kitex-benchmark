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

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"

	pb "github.com/cloudwego/kitex-benchmark/protobuf/rpcx/pb_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBRpcxClient(opt *runner.Options) runner.Client {
	cli := &pbRpcxClient{}
	cli.msg = string(opt.Body)
	cli.msgPool = &sync.Pool{
		New: func() interface{} {
			return &pb.RpcxMsg{}
		},
	}
	option := client.DefaultOption
	option.SerializeType = protocol.ProtoBuffer
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+opt.Address, "")
	cli.clipool = client.NewXClientPool(opt.PoolSize, "RpcxEcho", client.Failtry, client.RandomSelect, d, option)
	return cli
}

type pbRpcxClient struct {
	msg     string
	msgPool *sync.Pool
	clipool *client.XClientPool
}

func (cli *pbRpcxClient) Echo(msg string) (err error) {
	args := cli.msgPool.Get().(*pb.RpcxMsg)
	reply := cli.msgPool.Get().(*pb.RpcxMsg)
	defer cli.msgPool.Put(args)
	defer cli.msgPool.Put(reply)
	args.Msg = msg

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	xclient := cli.clipool.Get()
	err = xclient.Call(ctx, "EchoMsg", args, reply)
	return err
}
