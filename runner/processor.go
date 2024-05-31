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
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/kitex-benchmark/perf"
)

const (
	EchoAction        = "echo"
	EchoComplexAction = "echo_complex"
	BeginAction       = "begin"
	EndAction         = "end"
	SleepAction       = "sleep"
	ReportAction      = "report"
)

func ProcessRequest(recorder *perf.Recorder, action, msg string) (retAction, retMsg string) {
	switch action {
	case BeginAction:
		recorder.Begin()
	case EndAction:
		recorder.End()
		// report on server side
		recorder.Report()
		// send back server report to client
		return ReportAction, recorder.ReportString()
	case SleepAction:
		timeStr := strings.Split(msg, ",")[0]
		if n, err := strconv.Atoi(timeStr); err == nil {
			ms := time.Millisecond * time.Duration(n)
			if ms > 0 {
				time.Sleep(ms)
			}
		}
	default:
		// do business logic
	}

	return action, msg
}

func ProcessResponse(action, msg string) {
	switch action {
	case ReportAction:
		fmt.Print(msg)
	}
}
