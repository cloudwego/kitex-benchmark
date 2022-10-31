/*
 * Copyright 2022 CloudWeGo Authors
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
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/runner"
)

var requestData = []byte(`
{
    "msgMap":{
        "v1":{
            "id":123,
            "value":"hello"
        },
        "v2":{
            "id":321,
            "value":"world"
        }
    },
    "subMsgs":[
        {
            "id":123,
            "value":"hello"
        },
        {
            "id":321,
            "value":"world"
        }
    ],
    "msgSet":[
        {
            "id":123,
            "value":"hello",
            "subMessages":[
                {
                    "id":123,
                    "value":"hello"
                },
                {
                    "id":321,
                    "value":"world"
                }
            ]
        },
        {
            "id":321,
            "value":"world",
            "subMessages":[
                {
                    "id":321,
                    "value":"world"
                },
                {
                    "id":123,
                    "value":"hello"
                }
            ]
        }
    ],
    "flagMsg":{
        "id":123,
        "value":"hello",
        "subMessages":[
            {
                "id":123,
                "value":"hello"
            },
            {
                "id":321,
                "value":"world"
            }
        ]
    }
}
`)

func NewGenericHTTPClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造http 请求和返回类型的泛化调用
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericHTTPClient{}
	cli.client, err = genericclient.NewClient("test.echo.kitex", g,
		client.WithTransportProtocol(transport.TTHeader),
		client.WithHostPorts(opt.Address),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithLongConnection(
			connpool.IdleConfig{MaxIdlePerAddress: 1000, MaxIdleGlobal: 1000, MaxIdleTimeout: time.Minute}),
	)
	if err != nil {
		panic(err)
	}
	return cli
}

type genericHTTPClient struct {
	client genericclient.Client
}

func (cli *genericHTTPClient) Echo(action, msg string) error {
	ctx := context.Background()

	url := fmt.Sprintf("http://example.com/test/obj/%s", action)
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	httpRequest.Header.Set("msg", msg)

	// send the request
	customReq, err := generic.FromHTTPRequest(httpRequest)
	if err != nil {
		return err
	}

	reply, err := cli.client.GenericCall(ctx, "", customReq)
	if reply != nil {
		resp := reply.(*generic.HTTPResponse)
		runner.ProcessResponse(resp.Header.Get("action"), resp.Header.Get("msg"))
	}
	return err
}
