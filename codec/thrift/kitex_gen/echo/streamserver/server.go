// Code generated by Kitex v0.11.0. DO NOT EDIT.
package streamserver

import (
	"context"
	echo "github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex/pkg/streamx"
	"github.com/cloudwego/kitex/pkg/streamx/provider/ttstream"
	server "github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/streamxserver"
)

type Server interface {
	Echo(ctx context.Context, stream streamx.BidiStreamingServer[echo.Request, echo.Response]) error
}

func RegisterService(svr server.Server, handler Server, opts ...server.RegisterOption) error {
	sp, err := ttstream.NewServerProvider(svcInfo)
	if err != nil {
		return err
	}
	nopts := []server.RegisterOption{
		streamxserver.WithProvider(sp),
	}
	nopts = append(nopts, opts...)
	return svr.RegisterService(svcInfo, handler, nopts...)
}
