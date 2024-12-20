// Code generated by Kitex v0.11.3. DO NOT EDIT.

package streamserver

import (
	"context"
	"errors"
	echo "github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/streamxclient"
	"github.com/cloudwego/kitex/client/streamxclient/streamxcallopt"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/pkg/streamx"
	"github.com/cloudwego/kitex/server/streamxserver"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"Echo": kitex.NewMethodInfo(
		echoHandler,
		newStreamServerEchoArgs,
		newStreamServerEchoResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingBidirectional),
		kitex.WithMethodExtra("streamx", "true"),
	),
}

var (
	streamServerServiceInfo                = NewServiceInfo()
	streamServerServiceInfoForClient       = NewServiceInfoForClient()
	streamServerServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return streamServerServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return streamServerServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return streamServerServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(true, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "StreamServer"
	handlerType := (*echo.StreamServer)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "echo",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.11.3",
		Extra:           extra,
	}
	return svcInfo
}

func echoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	return streamxserver.InvokeBidiStreamHandler[echo.Request, echo.Response](
		ctx, arg.(streamx.StreamReqArgs), result.(streamx.StreamResArgs), func(ctx context.Context, stream streamx.BidiStreamingServer[echo.Request, echo.Response]) error {
			return handler.(StreamServer).Echo(ctx, stream)
		},
	)
}
func newStreamServerEchoArgs() interface{} {
	return echo.NewStreamServerEchoArgs()
}

func newStreamServerEchoResult() interface{} {
	return echo.NewStreamServerEchoResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Echo(ctx context.Context, callOptions ...streamxcallopt.CallOption) (
	context.Context, streamx.BidiStreamingClient[echo.Request, echo.Response], error,
) {
	return streamxclient.InvokeStream[echo.Request, echo.Response](
		ctx, p.c, kitex.StreamingBidirectional, "Echo", nil, nil, callOptions...)
}
