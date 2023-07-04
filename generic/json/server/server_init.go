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

package server

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
	"github.com/cloudwego/kitex-benchmark/generic/data"
	"github.com/cloudwego/kitex-benchmark/perf"
	"github.com/cloudwego/kitex-benchmark/runner"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

const (
	DefaultPort   = 8002
	DynamicGoPort = 8003
)

var recorder = perf.NewRecorder("GenericJSON@Server")

type GenericServerSmallImpl struct{}

func (s *GenericServerSmallImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "TestObj":
		req := request.(string)
		var rep echo.Response
		err = json.Unmarshal([]byte(req), &rep)
		if err != nil {
			return nil, kerrors.NewBizStatusError(500, err.Error())
		}
		resp := runner.ProcessRequest(recorder, rep.Action, rep.Msg)
		return data.GetJsonString(resp.Action, resp.Msg, data.Small), nil
	}
	return nil, kerrors.NewBizStatusError(404, "not found")
}

type GenericServerMediumImpl struct{}

func (s *GenericServerMediumImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "TestObj":
		req := request.(string)
		var rep echo.Response
		err = json.Unmarshal([]byte(req), &rep)
		if err != nil {
			return nil, kerrors.NewBizStatusError(500, err.Error())
		}
		resp := runner.ProcessRequest(recorder, rep.Action, rep.Msg)
		return data.GetJsonString(resp.Action, resp.Msg, data.Medium), nil
	}
	return nil, kerrors.NewBizStatusError(404, "not found")
}

type GenericServerLargeImpl struct{}

func (s *GenericServerLargeImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	switch method {
	case "TestObj":
		req := request.(string)
		var rep echo.Response
		err = json.Unmarshal([]byte(req), &rep)
		if err != nil {
			return nil, kerrors.NewBizStatusError(500, err.Error())
		}
		resp := runner.ProcessRequest(recorder, rep.Action, rep.Msg)
		return data.GetJsonString(resp.Action, resp.Msg, data.Large), nil
	}
	return nil, kerrors.NewBizStatusError(404, "not found")
}
