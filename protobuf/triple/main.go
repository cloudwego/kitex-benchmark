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
	api "github.com/cloudwego/kitex-benchmark/codec/protobuf/dubbo_gen"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
	"os"
)

type EchoProvider struct {
	api.UnimplementedEchoServer
}

var recorder = perf.NewRecorder("Dubbo-triple@Server")

func (s *EchoProvider) EchoHello(ctx context.Context, in *api.Request) (*api.Response, error) {
	resp := runner.ProcessRequest(recorder, in.Action, in.Msg)
	return &api.Response{Action: resp.Action, Msg: resp.Msg}, nil
}

func main() {
	os.Setenv("DUBBO_GO_CONFIG_PATH", "./protobuf/triple/dubbogo.yaml")
	config.SetProviderService(&EchoProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
