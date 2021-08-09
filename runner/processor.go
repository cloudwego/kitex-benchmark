package runner

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/kitex-benchmark/perf"
)

const (
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
		if n, err := strconv.Atoi(msg); err != nil {
			sleepTime := time.Millisecond * time.Duration(n)
			if sleepTime > 0 {
				time.Sleep(sleepTime)
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
		fmt.Printf(msg)
	}
}
