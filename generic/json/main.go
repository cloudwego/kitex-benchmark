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
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
)

const (
	port = 8002
)

var recorder = perf.NewRecorder("GenericJSON@Server")

var (
	responseData = `
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
}`
	actionidx, msgidx int
)

func init() {
	actionidx = strings.Index(responseData, `"action":""`) + len(`"action":""`) - 1
	msgidx = strings.Index(responseData, `"msg":""`) + len(`"msg":""`) - 1
}

func GetJsonString(action, msg string) string {
	return responseData[:actionidx] + action + responseData[actionidx:msgidx] + msg + responseData[msgidx:]
}

type GenericServerImpl struct{}

func (s *GenericServerImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "EchoComplex":
		req := request.(string)
		var rep echo.Response
		err = json.Unmarshal([]byte(req), &rep)
		if err != nil {
			return nil, kerrors.NewBizStatusError(500, err.Error())
		}
		resp := runner.ProcessRequest(recorder, rep.Action, rep.Msg)
		return GetJsonString(resp.Action, resp.Msg), nil
	}
	return nil, kerrors.NewBizStatusError(404, "not found")
}

func main() {
	// start pprof server
	go func() {
		perf.ServeMonitor(fmt.Sprintf(":%d", port+10000))
	}()

	// CurDir: ./scripts
	p, err := generic.NewThriftFileProvider("./codec/thrift/echo.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	address := &net.UnixAddr{Net: "tcp", Name: fmt.Sprintf(":%d", port)}
	svr := genericserver.NewServer(new(GenericServerImpl), g, server.WithServiceAddr(address), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
