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
