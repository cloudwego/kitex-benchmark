// Code generated by Kitex v0.11.3. DO NOT EDIT.

package echoserver

import (
	"context"
	echo "github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Echo(ctx context.Context, req *echo.Request, callOptions ...callopt.Option) (r *echo.Response, err error)
	EchoComplex(ctx context.Context, req *echo.ComplexRequest, callOptions ...callopt.Option) (r *echo.ComplexResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kEchoServerClient{
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

type kEchoServerClient struct {
	*kClient
}

func (p *kEchoServerClient) Echo(ctx context.Context, req *echo.Request, callOptions ...callopt.Option) (r *echo.Response, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Echo(ctx, req)
}

func (p *kEchoServerClient) EchoComplex(ctx context.Context, req *echo.ComplexRequest, callOptions ...callopt.Option) (r *echo.ComplexResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.EchoComplex(ctx, req)
}
