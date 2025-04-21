/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex-tests/pkg/utils"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/runner"
)

// main is use for routing.
func main() {
	runner.Main("KITEX-HTTP", NewThriftKitexClient)
}

func NewThriftKitexClient(opt *runner.Options) runner.Client {
	cli := &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
			MaxIdleConns:        1000,
			IdleConnTimeout:     time.Minute,
		},
	}
	return NewKitexClient(cli, opt.Address)
}

type kitexClient struct {
	client  *http.Client
	addr    string
	request []byte
}

func NewKitexClient(client *http.Client, addr string) *kitexClient {
	c := new(kitexClient)
	c.client = client
	c.addr = addr
	return c
}

func (cli *kitexClient) Send(method, action, msg string) error {
	req := createComplexRequest(action, msg)
	buf, _ := sonic.Marshal(req)
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/EchoService/EchoComplex", cli.addr), bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	response, err := cli.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	resp := &echo.ComplexResponse{}
	rresp := &Response{Data: resp}
	err = sonic.Unmarshal(body, rresp)
	if err != nil {
		return err
	}
	runner.ProcessResponse(resp.Action, resp.Msg)
	return nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var echoComplexReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.ComplexRequest{}
	},
}

func createComplexRequest(action, msg string) *echo.ComplexRequest {
	req := echoComplexReqPool.Get().(*echo.ComplexRequest)
	defer echoComplexReqPool.Put(req)

	id := int64(fastrand.Int31n(100))
	smallSubMsg := &echo.SubMessage{
		Id:    &id,
		Value: ptr(utils.RandomString(10)),
	}
	subMsg1K := &echo.SubMessage{
		Id:    &id,
		Value: ptr(utils.RandomString(1024)),
	}

	subMsgList2Items := []*echo.SubMessage{smallSubMsg, smallSubMsg}

	message := &echo.Message{
		Id:          &id,
		Value:       ptr(utils.RandomString(1024)),
		SubMessages: subMsgList2Items,
	}

	msgMap := make(map[string]*echo.SubMessage)
	for i := 0; i < 5; i++ {
		msgMap[strconv.Itoa(i)] = subMsg1K
	}

	subMsgList100Items := make([]*echo.SubMessage, 100)
	for i := 0; i < len(subMsgList100Items); i++ {
		subMsgList100Items[i] = smallSubMsg
	}

	req.Action = action
	req.Msg = msg
	req.MsgMap = msgMap
	req.SubMsgs = subMsgList100Items
	req.MsgSet = []*echo.Message{message}
	req.FlagMsg = message

	return req
}

func ptr[T any](v T) *T { return &v }
