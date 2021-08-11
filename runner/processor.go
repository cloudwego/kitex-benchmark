package runner

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/kitex-benchmark/perf"
)

const (
	EchoAction   = "echo"
	BeginAction  = "begin"
	EndAction    = "end"
	SleepAction  = "sleep"
	ReportAction = "report"
)

func ProcessRequest(recorder *perf.Recorder, action, msg string) *Response {
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
	default:
		// do business logic
	}

	return &Response{
		Action: action,
		Msg:    msg,
	}
}

func ProcessResponse(action, msg string) {
	switch action {
	case ReportAction:
		fmt.Print(msg)
	}
}
