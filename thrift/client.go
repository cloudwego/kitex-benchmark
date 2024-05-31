package thrift

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo/echoserver"
	"github.com/cloudwego/kitex-benchmark/runner"
)

type kitexClient struct {
	client echoserver.Client
}

func NewKitexClient(client echoserver.Client) *kitexClient {
	c := new(kitexClient)
	c.client = client
	return c
}

func (cli *kitexClient) Send(method, action, msg string) error {
	switch strings.ToLower(method) {
	case "echo":
		return cli.echo(action, msg)
	case "echocomplex":
		return cli.echoComplex(action, msg)
	default:
		return fmt.Errorf("unknow method: %s", method)
	}
}

var echoReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.Request{}
	},
}

func (cli *kitexClient) echo(action, msg string) error {
	ctx := context.Background()
	req := echoReqPool.Get().(*echo.Request)
	defer echoReqPool.Put(req)

	req.Action = action
	req.Msg = msg

	reply, err := cli.client.Echo(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}

var echoComplexReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.ComplexRequest{}
	},
}

func (cli *kitexClient) echoComplex(action, msg string) error {
	ctx := context.Background()
	req := echoComplexReqPool.Get().(*echo.ComplexRequest)
	defer echoComplexReqPool.Put(req)

	// 复杂结构体下，我们需要把 msg string 分拆到 complex request 中，保证整体包大小没有太大变化下，提高字段复杂度
	const complexity = 16
	msgSize := len(msg)
	req.Action = action

	req.MsgMap = make(map[string]*echo.SubMessage, complexity)
	content := msg[msgSize/4*0 : msgSize/4*1]
	for idx, str := range splitString(content, complexity) {
		id := int64(idx)
		req.MsgMap[strconv.Itoa(idx)] = &echo.SubMessage{
			Id:    &id,
			Value: &str,
		}
	}

	req.SubMsgs = make([]*echo.SubMessage, complexity)
	content = msg[msgSize/4*1 : msgSize/4*2]
	for idx, str := range splitString(content, complexity) {
		id := int64(idx)
		req.SubMsgs[idx] = &echo.SubMessage{
			Id:    &id,
			Value: &str,
		}
	}

	req.MsgSet = make([]*echo.Message, complexity)
	content = msg[msgSize/4*2 : msgSize/4*3]
	for idx, str := range splitString(content, complexity) {
		id := int64(idx)
		req.MsgSet[idx] = &echo.Message{
			Id: &id,
			SubMessages: []*echo.SubMessage{
				{
					Id:    &id,
					Value: &str,
				},
			},
		}
	}

	req.FlagMsg = new(echo.Message)
	content = msg[msgSize/4*3 : msgSize/4*4]
	req.FlagMsg = &echo.Message{
		Value: &content,
	}

	reply, err := cli.client.EchoComplex(ctx, req)
	if reply != nil {
		runner.ProcessResponse(reply.Action, reply.Msg)
	}
	return err
}

func splitString(str string, n int) []string {
	ret := make([]string, n)
	if n < 0 || len(str) < n {
		return ret
	}
	single := len(str) / n
	for i := 0; i < n; i++ {
		ret[i] = str[i*single : (i+1)*single]
	}
	return ret
}
