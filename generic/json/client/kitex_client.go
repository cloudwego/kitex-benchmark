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
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"

	"github.com/cloudwego/kitex-benchmark/runner"
)

var (
	requestData = `
{
    "action":"",
    "msg":"",
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
`
	actionidx, msgidx int
)

func init() {
	actionidx = strings.Index(requestData, `"action":""`) + len(`"action":""`) - 1
	msgidx = strings.Index(requestData, `"msg":""`) + len(`"msg":""`) - 1
}

func GetJsonString(action, msg string) string {
	return requestData[:actionidx] + action + requestData[actionidx:msgidx] + msg + requestData[msgidx:]
}

func NewGenericJSONClient(opt *runner.Options) runner.Client {
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	// 构造json 请求和返回类型的泛化调用
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	cli := &genericJSONClient{}
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

type genericJSONClient struct {
	client genericclient.Client
}

func (cli *genericJSONClient) Echo(action, msg string) error {
	ctx := context.Background()

	reply, err := cli.client.GenericCall(ctx, "EchoComplex", GetJsonString(action, msg))
	if reply != nil {
		repl := reply.(string)
		var rep echo.Request
		err = json.Unmarshal([]byte(repl), &rep)
		if err != nil {
			return err
		}
		runner.ProcessResponse(rep.Action, rep.Msg)
	}
	return err
}
