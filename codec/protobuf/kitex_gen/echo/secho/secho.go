// Code generated by Kitex v0.1.4. DO NOT EDIT.

package secho

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex-benchmark/codec/protobuf/kitex_gen/echo"
	"github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/pkg/streaming"
	"google.golang.org/protobuf/proto"
)

func serviceInfo() *kitex.ServiceInfo {
	return sEchoServiceInfo
}

var sEchoServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "SEcho"
	handlerType := (*echo.SEcho)(nil)
	methods := map[string]kitex.MethodInfo{
		"echo": kitex.NewMethodInfo(echoHandler, newEchoArgs, newEchoResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "protobuf",
	}
	extra["streaming"] = true
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.1.4",
		Extra:           extra,
	}
	return svcInfo
}

func echoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	st := arg.(*streaming.Args).Stream
	stream := &sEchoechoServer{st}
	return handler.(echo.SEcho).Echo(stream)
}

type sEchoechoClient struct {
	streaming.Stream
}

func (x *sEchoechoClient) Send(m *echo.SRequest) error {
	return x.Stream.SendMsg(m)
}
func (x *sEchoechoClient) Recv() (*echo.SResponse, error) {
	m := new(echo.SResponse)
	return m, x.Stream.RecvMsg(m)
}

type sEchoechoServer struct {
	streaming.Stream
}

func (x *sEchoechoServer) Send(m *echo.SResponse) error {
	return x.Stream.SendMsg(m)
}

func (x *sEchoechoServer) Recv() (*echo.SRequest, error) {
	m := new(echo.SRequest)
	return m, x.Stream.RecvMsg(m)
}

func newEchoArgs() interface{} {
	return &EchoArgs{}
}

func newEchoResult() interface{} {
	return &EchoResult{}
}

type EchoArgs struct {
	Req *echo.SRequest
}

func (p *EchoArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in EchoArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *EchoArgs) Unmarshal(in []byte) error {
	msg := new(echo.SRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var EchoArgs_Req_DEFAULT *echo.SRequest

func (p *EchoArgs) GetReq() *echo.SRequest {
	if !p.IsSetReq() {
		return EchoArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *EchoArgs) IsSetReq() bool {
	return p.Req != nil
}

type EchoResult struct {
	Success *echo.SResponse
}

var EchoResult_Success_DEFAULT *echo.SResponse

func (p *EchoResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in EchoResult")
	}
	return proto.Marshal(p.Success)
}

func (p *EchoResult) Unmarshal(in []byte) error {
	msg := new(echo.SResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *EchoResult) GetSuccess() *echo.SResponse {
	if !p.IsSetSuccess() {
		return EchoResult_Success_DEFAULT
	}
	return p.Success
}

func (p *EchoResult) SetSuccess(x interface{}) {
	p.Success = x.(*echo.SResponse)
}

func (p *EchoResult) IsSetSuccess() bool {
	return p.Success != nil
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Echo(ctx context.Context) (SEcho_echoClient, error) {
	streamClient, ok := p.c.(client.Streaming)
	if !ok {
		return nil, fmt.Errorf("client not support streaming")
	}
	res := new(streaming.Result)
	err := streamClient.Stream(ctx, "echo", nil, res)
	if err != nil {
		return nil, err
	}
	stream := &sEchoechoClient{res.Stream}
	return stream, nil
}