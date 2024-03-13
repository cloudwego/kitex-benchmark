package main

import (
	"context"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/protobuf/hessian/client/pkg"
	"github.com/cloudwego/kitex-benchmark/runner"
)

type EchoProvider struct {
	CommonEchoProvider
}
type CommonEchoProvider struct {
}

var recorder = perf.NewRecorder("Dubbo@Server")

func (c *CommonEchoProvider) Echo(ctx context.Context, in *pkg.Message) (*pkg.Message, error) {
	resp := runner.ProcessRequest(recorder, in.Action, in.Msg)
	return &pkg.Message{Action: resp.Action, Msg: resp.Msg}, nil
}
