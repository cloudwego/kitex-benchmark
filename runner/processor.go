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

package runner

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/kitex-benchmark/perf"
)

const (
	EchoAction    = "echo"
	BeginAction   = "begin"
	EndAction     = "end"
	SleepAction   = "sleep"
	ComplexAction = "complex"
	ReportAction  = "report"
)

// complexProcess will consume CPU and Memory, and create more GC work
func complexProcess(concurrency, size int) []byte {
	var wg sync.WaitGroup
	heap := make([][]byte, concurrency)
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			heap[c] = make([]byte, size)
			for i := 0; i < len(heap[c]); i++ {
				heap[c][i] = byte('a' + ((i + c) % 26))
			}
		}(c)
	}
	wg.Wait()
	return heap[concurrency/2]
}

func ProcessRequest(recorder *perf.Recorder, action, msg string) *Response {
	respMsg := msg
	switch action {
	case BeginAction:
		recorder.Begin()
	case EndAction:
		recorder.End()
		// report on server side
		recorder.Report()
		// send back server report to client
		return &Response{
			Action: ReportAction,
			Msg:    recorder.ReportString(),
		}
	case SleepAction:
		timeStr := strings.Split(msg, ",")[0]
		if n, err := strconv.Atoi(timeStr); err == nil {
			ms := time.Millisecond * time.Duration(n)
			if ms > 0 {
				time.Sleep(ms)
			}
		}
	case ComplexAction:
		respMsg = string(complexProcess(runtime.GOMAXPROCS(0), len(msg)))
	default:
		// do business logic
	}

	return &Response{
		Action: action,
		Msg:    respMsg,
	}
}

func ProcessResponse(action, msg string) {
	switch action {
	case ReportAction:
		fmt.Print(msg)
	}
}
