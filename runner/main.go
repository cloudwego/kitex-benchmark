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
	flag.StringVar(&address, "addr", "", "client call address")
	flag.IntVar(&echoSize, "b", 1024, "echo size once")
	flag.IntVar(&concurrent, "c", 1, "call concurrent")
	flag.Int64Var(&total, "n", 1, "call total nums")
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
	recorder := perf.NewRecorder(fmt.Sprintf("%s@Client", name))
	recorder.Begin()

	// === benching ===
	r.Run(name, handler, concurrent, total, echoSize, sleepTime)

	// == ending ===
	recorder.End()
	if err := cli.Echo(EndAction, ""); err != nil {
		log.Fatalf("ending server failed: %v", err)
	}

	// === reporting ===
	recorder.Report() // report client
	fmt.Printf("\n\n")
}
