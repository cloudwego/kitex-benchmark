// Code generated by Kitex v0.11.0. DO NOT EDIT.

package secho

import (
	"context"
	echo "github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	transport "github.com/cloudwego/kitex/transport"
	"github.com/cloudwego/kitex/client/streamclient"
	"github.com/cloudwego/kitex/client/callopt/streamcall"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Echo(ctx context.Context, callOptions ...callopt.Option) (stream SEcho_echoClient, err error)
}

// StreamClient is designed to provide Interface for Streaming APIs.
type StreamClient interface {
	Echo(ctx context.Context, callOptions ...streamcall.Option) (stream SEcho_echoClient, err error)
}

type SEcho_echoClient interface {
	streaming.Stream
	Send(*echo.Request) error
	Recv() (*echo.Response, error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, client.WithTransportProtocol(transport.GRPC))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kSEchoClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kSEchoClient struct {
	*kClient
}

func (p *kSEchoClient) Echo(ctx context.Context, callOptions ...callopt.Option) (stream SEcho_echoClient, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Echo(ctx)
}

// NewStreamClient creates a stream client for the service's streaming APIs defined in IDL.
func NewStreamClient(destService string, opts ...streamclient.Option) (StreamClient, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))
	options = append(options, client.WithTransportProtocol(transport.GRPC))
	options = append(options, streamclient.GetClientOptions(opts)...)

	kc, err := client.NewClient(serviceInfoForStreamClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kSEchoStreamClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewStreamClient creates a stream client for the service's streaming APIs defined in IDL.
// It panics if any error occurs.
func MustNewStreamClient(destService string, opts ...streamclient.Option) StreamClient {
	kc, err := NewStreamClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kSEchoStreamClient struct {
	*kClient
}

func (p *kSEchoStreamClient) Echo(ctx context.Context, callOptions ...streamcall.Option) (stream SEcho_echoClient, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, streamcall.GetCallOptions(callOptions))
	return p.kClient.Echo(ctx)
}
