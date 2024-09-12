// Code generated by Kitex v0.11.0. DO NOT EDIT.

package streamserver

import (
	"context"
	echo "github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/streamxclient"
	"github.com/cloudwego/kitex/client/streamxclient/streamxcallopt"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/pkg/streamx"
	"github.com/cloudwego/kitex/pkg/streamx/provider/ttstream"
)

type Client interface {
	Echo(ctx context.Context, callOptions ...streamxcallopt.CallOption) (stream streamx.BidiStreamingClient[ttstream.Header, ttstream.Trailer, echo.Request, echo.Response], err error)
}

func NewClient(destService string, opts ...streamxclient.Option) (Client, error) {
	var options []streamxclient.Option
	options = append(options, streamxclient.WithDestService(destService))
	options = append(options, opts...)
	cp, err := ttstream.NewClientProvider(svcInfo)
	if err != nil {
		return nil, err
	}
	options = append(options, streamxclient.WithProvider(cp))
	cli, err := streamxclient.NewClient(svcInfo, options...)
	if err != nil {
		return nil, err
	}
	kc := &kClient{streamer: cli, caller: cli.(client.Client)}
	return kc, nil
}

var _ Client = (*kClient)(nil)

type kClient struct {
	caller   client.Client
	streamer streamxclient.Client
}

func (c *kClient) Echo(ctx context.Context, callOptions ...streamxcallopt.CallOption) (stream streamx.BidiStreamingClient[ttstream.Header, ttstream.Trailer, echo.Request, echo.Response], err error) {
	return streamxclient.InvokeStream[ttstream.Header, ttstream.Trailer, echo.Request, echo.Response](
		ctx, c.streamer, serviceinfo.StreamingBidirectional, "Echo", nil, nil, callOptions...)
}
