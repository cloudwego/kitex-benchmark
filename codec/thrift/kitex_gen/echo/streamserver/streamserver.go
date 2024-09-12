// Code generated by Kitex v0.11.0. DO NOT EDIT.
package streamserver

import (
	"context"
	echo "github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/pkg/streamx"
	ttheader "github.com/cloudwego/kitex/pkg/streamx/provider/ttstream"
	"github.com/cloudwego/kitex/server/streamxserver"
)

var svcInfo = &serviceinfo.ServiceInfo{
	ServiceName: "StreamServer",
	Methods: map[string]serviceinfo.MethodInfo{
		"Echo": serviceinfo.NewMethodInfo(
			func(ctx context.Context, handler, reqArgs, resArgs interface{}) error {
				return streamxserver.InvokeStream[ttheader.Header, ttheader.Trailer, echo.Request, echo.Response](
					ctx, serviceinfo.StreamingBidirectional, handler.(streamx.StreamHandler), reqArgs.(streamx.StreamReqArgs), resArgs.(streamx.StreamResArgs))
			},
			nil,
			nil,
			false,
			serviceinfo.WithStreamingMode(serviceinfo.StreamingBidirectional),
		),
	},
	Extra: map[string]interface{}{
		"streaming": true,
	},
}

func ServiceInfo() *serviceinfo.ServiceInfo {
	return svcInfo
}
