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

package data

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cloudwego/kitex-benchmark/codec/thrift/kitex_gen/echo"
)

var (
	SmallData, MediumData, LargeData                              string
	actionSidx, msgSidx, actionMidx, msgMidx, actionLidx, msgLidx int
)

type Size int

const (
	Small Size = iota
	Medium
	Large
)

func init() {
	//TODO: omit action and msg for http?
	SmallData = getReqValue(1)
	MediumData = getReqValue(10)
	LargeData = getReqValue(100)
	actionSidx = strings.Index(SmallData, `"action":""`) + len(`"action":""`) - 1
	actionMidx = strings.Index(MediumData, `"action":""`) + len(`"action":""`) - 1
	actionLidx = strings.Index(LargeData, `"action":""`) + len(`"action":""`) - 1
	msgSidx = strings.Index(SmallData, `"msg":""`) + len(`"msg":""`) - 1
	msgMidx = strings.Index(MediumData, `"msg":""`) + len(`"msg":""`) - 1
	msgLidx = strings.Index(LargeData, `"msg":""`) + len(`"msg":""`) - 1
}

func GetJsonString(action, msg string, size Size) string {
	switch size {
	case Small:
		return SmallData[:actionSidx] + action + SmallData[actionSidx:msgSidx] + msg + SmallData[msgSidx:]
	case Medium:
		return MediumData[:actionMidx] + action + MediumData[actionMidx:msgMidx] + msg + MediumData[msgMidx:]
	case Large:
		return LargeData[:actionLidx] + action + LargeData[actionLidx:msgLidx] + msg + LargeData[msgLidx:]
	}
	return ""
}

func getReqValue(size int) string {
	req := &echo.ObjReq{
		Action:  "",
		Msg:     "",
		MsgMap:  map[string]*echo.SubMessage{},
		SubMsgs: []*echo.SubMessage{},
		MsgSet:  []*echo.Message{},
		FlagMsg: &echo.Message{},
	}

	for i := 0; i < size; i++ {
		req.MsgMap[strconv.Itoa(i)] = getSubMessage(int64(i))
		req.SubMsgs = append(req.SubMsgs, getSubMessage(int64(i)))
		req.MsgSet = append(req.MsgSet, getMessage(int64(i)))
		req.FlagMsg = getMessage(int64(i))
	}

	data, _ := json.Marshal(req)
	return string(data)
}

func getSubMessage(i int64) *echo.SubMessage {
	value := "hello"
	return &echo.SubMessage{
		Id:    &i,
		Value: &value,
	}
}

func getMessage(i int64) *echo.Message {
	value := "hello"
	ret := &echo.Message{
		Id:          &i,
		Value:       &value,
		SubMessages: []*echo.SubMessage{},
	}
	ret.SubMessages = append(ret.SubMessages, getSubMessage(1))
	ret.SubMessages = append(ret.SubMessages, getSubMessage(2))
	return ret
}
