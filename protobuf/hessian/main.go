package main

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/cloudwego/kitex-benchmark/protobuf/hessian/client/pkg"
	"os"
)

func main() {
	os.Setenv("DUBBO_GO_CONFIG_PATH", "./protobuf/hessian/dubbogo.yaml")
	hessian.RegisterPOJO(&pkg.Message{})
	config.SetProviderService(&EchoProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
