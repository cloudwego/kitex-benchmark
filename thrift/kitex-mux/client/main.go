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
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex-benchmark/thrift"
)

// main is use for routing.
func main() {
	runner.Main("KITEX-MUX", newThriftKitexClient)
}

func newThriftKitexClient(opt *runner.Options) runner.Client {
	cli := echoserver.MustNewClient("test.echo.kitex",
		client.WithTransportProtocol(transport.Framed),
		client.WithHostPorts(opt.Address),
		client.WithMuxConnection(2))
	return thrift.NewKitexClient(cli)
}
