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

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-server/conf/dubbogo.yaml
func main() {
	os.Setenv("DUBBO_GO_CONFIG_PATH", "./protobuf/triple/dubbogo.yaml")
	config.SetProviderService(&EchoProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
