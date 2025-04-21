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
	"sync"
	"time"

	"github.com/bytedance/sonic"

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

var echoReqPool = sync.Pool{
	New: func() interface{} {
		return &echo.Request{}
	},
}

func (cli *kitexClient) Send(method, action, msg string) error {
	req := echoReqPool.Get().(*echo.Request)
	defer echoReqPool.Put(req)

	req.Action = action
	req.Msg = msg
	buf, _ := sonic.Marshal(req)
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/EchoServer/Echo", cli.addr), bytes.NewBuffer(buf))
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
	resp := &echo.Response{}
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
