/*
 * Copyright 2022 CloudWeGo Authors
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
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"

	kserver "github.com/cloudwego/kitex-benchmark/generic/json/server"
	"github.com/cloudwego/kitex-benchmark/perf"
)

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", kserver.DefaultPort+10000))
	}()

	// CurDir: ./scripts
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", kserver.DefaultPort)}
	svr := genericserver.NewServer(new(kserver.GenericServerSmallImpl), g, server.WithServiceAddr(address), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
