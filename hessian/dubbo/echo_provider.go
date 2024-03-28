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
	"github.com/cloudwego/kitex-benchmark/hessian/dubbo/client/pkg"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

type EchoProvider struct {
	CommonEchoProvider
}
type CommonEchoProvider struct {
}

var recorder = perf.NewRecorder("Hessian@Server")

func (c *CommonEchoProvider) Echo(_ context.Context, in *pkg.Message) (*pkg.Message, error) {
	resp := runner.ProcessRequest(recorder, in.Action, in.Msg)
	return &pkg.Message{Action: resp.Action, Msg: resp.Msg}, nil
}
