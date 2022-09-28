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
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/cloudwego/kitex-benchmark/perf"
)

var (
	address    string
	echoSize   int
	total      int64
	concurrent int
	poolSize   int
	sleepTime  int
)

type Options struct {
	Address  string
	Body     []byte
	PoolSize int
}

type ClientNewer func(opt *Options) Client

type Client interface {
	Echo(action, msg string) (err error)
}

type Response struct {
	Action string
	Msg    string
}

func initFlags() {
	flag.StringVar(&address, "addr", "127.0.0.1:8000", "client call address")
	flag.IntVar(&echoSize, "b", 1024, "echo size once")
	flag.IntVar(&concurrent, "c", 100, "call concurrent")
	flag.Int64Var(&total, "n", 1024*100, "call total nums")
	flag.IntVar(&poolSize, "pool", 10, "conn poll size")
	flag.IntVar(&sleepTime, "sleep", 0, "sleep time for every request handler")
	flag.Parse()
}

func Main(name string, newer ClientNewer) {
	initFlags()

	// start pprof server
	go func() {
		err := perf.ServeMonitor(":18888")
		if err != nil {
			fmt.Printf("perf monitor server start failed: %v\n", err)
		} else {
			fmt.Printf("perf monitor server start success\n")
		}
	}()

	r := NewRunner()

	opt := &Options{
		Address:  address,
		PoolSize: poolSize,
	}
	cli := newer(opt)
	payload := string(make([]byte, echoSize))
	action := EchoAction
	if sleepTime > 0 {
		action = SleepAction
		st := strconv.Itoa(sleepTime)
		payload = fmt.Sprintf("%s,%s", st, payload[len(st)+1:])
	}
	handler := func() error { return cli.Echo(action, payload) }

	// === warming ===
	r.Warmup(handler, concurrent, 100*1000)

	// === beginning ===
	if err := cli.Echo(BeginAction, ""); err != nil {
		log.Fatalf("beginning server failed: %v", err)
	}
	// periodical test
	periodicalTest(name, handler, concurrent, total)

	if err := cli.Echo(EndAction, ""); err != nil {
		log.Fatalf("ending server failed: %v", err)
	}

	fmt.Printf("\n\n")
}

func periodicalTest(name string, handler func() error, concurrency int, totalReq int64) {
	fmt.Printf("Start periodical test\n")
	var (
		round = 10
		reqNum = totalReq / int64(round)
		maxIdleTimeout = 3 * time.Second
	)

	var costs []int64
	var totalns, actualTotal, failed int64

	for i := 0; i < round; i++ {
		roundRunner := NewRunner()
		t, counter := roundRunner.Run(name, handler, concurrency, reqNum, echoSize, 0, false)
		costs = append(costs, counter.Costs()...)
		totalns += t
		actualTotal += counter.Total
		failed += counter.Failed

		time.Sleep(2 * maxIdleTimeout)
	}

	res := report(name, totalns, actualTotal, failed, costs, concurrency)
	fmt.Printf("%s\n", res)
}

func report(name string, totalns, total, failed int64, costs []int64, concurrency int) string {
	sec := int64(time.Second)
	var tps float64
	if totalns < sec {
		tps = float64(total*sec) / float64(totalns)
	} else {
		tps = float64(total) / (float64(totalns) / float64(sec))
	}
	fcosts := make([]float64, len(costs))
	for i := range costs {
		fcosts[i] = float64(costs[i])
	}
	tp99, _ := stats.Percentile(fcosts, 99)
	tp999, _ := stats.Percentile(fcosts, 99.9)
	var result string
	if tp999/1000 < 1 {
		result = fmt.Sprintf("[%s]: TPS: %.2f, TP99: %.2fus, TP999: %.2fus (b=%d Byte, c=%d, n=%d, failed=%d)",
			name, tps, tp99/1000, tp999/1000, echoSize, concurrency, total, failed)
	} else {
		result = fmt.Sprintf("[%s]: TPS: %.2f, TP99: %.2fms, TP999: %.2fms (b=%d Byte, c=%d, n=%d, failed=%d)",
			name, tps, tp99/1000000, tp999/1000000, echoSize, concurrency, total, failed)
	}
	return result
}