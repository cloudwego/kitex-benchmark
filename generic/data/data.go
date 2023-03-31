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
	SmallReq, MediumReq, LargeReq                                 *echo.ObjReq
	SmallMap, MediumMap, LargeMap                                 map[string]interface{}
	SmallString, MediumString, LargeString                        string
	actionSidx, msgSidx, actionMidx, msgMidx, actionLidx, msgLidx int
)

type Size int

const (
	Small  Size = 6
	Medium Size = 33
	Large  Size = 67
)

func init() {
	// data size: small 1027B, medium 5035B, large 10101B
	SmallReq = getReqValue(int(Small))
	MediumReq = getReqValue(int(Medium))
	LargeReq = getReqValue(int(Large))
	SmallString = reqToString(SmallReq)
	MediumString = reqToString(MediumReq)
	LargeString = reqToString(LargeReq)
	SmallMap = getReqMap(int(Small))
	MediumMap = getReqMap(int(Medium))
	LargeMap = getReqMap(int(Large))
	actionSidx = strings.Index(SmallString, `"action":""`) + len(`"action":""`) - 1
	actionMidx = strings.Index(MediumString, `"action":""`) + len(`"action":""`) - 1
	actionLidx = strings.Index(LargeString, `"action":""`) + len(`"action":""`) - 1
	msgSidx = strings.Index(SmallString, `"msg":""`) + len(`"msg":""`) - 1
	msgMidx = strings.Index(MediumString, `"msg":""`) + len(`"msg":""`) - 1
	msgLidx = strings.Index(LargeString, `"msg":""`) + len(`"msg":""`) - 1
}

func GetJsonString(action, msg string, size Size) string {
	switch size {
	case Small:
		return SmallString[:actionSidx] + action + SmallString[actionSidx:msgSidx] + msg + SmallString[msgSidx:]
	case Medium:
		return MediumString[:actionMidx] + action + MediumString[actionMidx:msgMidx] + msg + MediumString[msgMidx:]
	case Large:
		return LargeString[:actionLidx] + action + LargeString[actionLidx:msgLidx] + msg + LargeString[msgLidx:]
	}
	return ""
}

func reqToString(req *echo.ObjReq) string {
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func getReqValue(size int) *echo.ObjReq {
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

	return req
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

func getReqMap(size int) map[string]interface{} {
	var msgMap map[interface{}]interface{}
	var subMsgs []interface{}
	var msgSet []interface{}
	var flagMsg map[string]interface{}
	for i := 0; i < size; i++ {
		msgMap[strconv.Itoa(i)] = getSubMessageMap(int64(i))
		subMsgs = append(subMsgs, getSubMessageMap(int64(i)))
		msgSet = append(msgSet, getMessageMap(int64(i)))
		flagMsg = getMessageMap(int64(i))
	}
	return map[string]interface{}{
		"msgMap":  msgMap,
		"subMsgs": subMsgs,
		"msgSet":  msgSet,
		"flagMsg": flagMsg,
	}
}

func getSubMessageMap(i int64) map[string]interface{} {
	return map[string]interface{}{
		"id":    i,
		"value": "hello",
	}
}

func getMessageMap(i int64) map[string]interface{} {
	return map[string]interface{}{
		"id":          i,
		"value":       "hello",
		"subMessages": []interface{}{getSubMessageMap(1), getSubMessageMap(2)},
	}
}
