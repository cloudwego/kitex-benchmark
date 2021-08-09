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

	gogo "github.com/cloudwego/kitex-benchmark/codec/protobuf/gogo_gen"
	"github.com/cloudwego/kitex-benchmark/runner"
)

func NewPBRpcxClient(opt *runner.Options) runner.Client {
	cli := &pbRpcxClient{}
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

	option := client.DefaultOption
	option.SerializeType = protocol.ProtoBuffer
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+opt.Address, "")
	cli.clipool = client.NewXClientPool(opt.PoolSize, "Echo", client.Failtry, client.RandomSelect, d, option)
	return cli
}

type pbRpcxClient struct {
	msg      string
	reqPool  *sync.Pool
	respPool *sync.Pool
	clipool  *client.XClientPool
}

func (cli *pbRpcxClient) Echo(action, msg string) (err error) {
	args := cli.reqPool.Get().(*gogo.Request)
	reply := cli.respPool.Get().(*gogo.Response)
	defer cli.reqPool.Put(args)
	defer cli.respPool.Put(reply)

	args.Action = action
	args.Msg = msg

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	xclient := cli.clipool.Get()
	err = xclient.Call(ctx, "Echo", args, reply)

	runner.ProcessResponse(reply.Action, reply.Msg)
	return err
}
