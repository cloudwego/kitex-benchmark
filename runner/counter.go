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
	"sync/atomic"
	"time"

	"github.com/montanaflynn/stats"
)

// 计数器
type Counter struct {
	Total  int64   // 总调用次数(limiter)
	Failed int64   // 失败次数
	costs  []int64 // 耗时统计
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Reset(total int64) {
	atomic.StoreInt64(&c.Total, 0)
	atomic.StoreInt64(&c.Failed, 0)
	c.costs = make([]int64, total)
}

func (c *Counter) AddRecord(idx int64, err error, cost int64) {
	c.costs[idx] = cost
	if err != nil {
		atomic.AddInt64(&c.Failed, 1)
	}
}

// idx < 0 表示用尽
func (c *Counter) Idx() (idx int64) {
	return atomic.AddInt64(&c.Total, 1) - 1
}

func (c *Counter) Report(title string, totalns int64, concurrent int, total int64, echoSize int) error {
	ms, sec := int64(time.Millisecond), int64(time.Second)
	logInfo("[%s]: finish benching [%s], took %d ms for %d requests", title, time.Now().String(), totalns/ms, c.Total)
	logInfo("[%s]: requests total: %d, failed: %d", title, c.Total, c.Failed)

	var tps float64
	if totalns < sec {
		tps = float64(c.Total*sec) / float64(totalns)
	} else {
		tps = float64(c.Total) / (float64(totalns) / float64(sec))
	}

	costs := make([]float64, len(c.costs))
	for i := range c.costs {
		costs[i] = float64(c.costs[i])
	}
	tp99, _ := stats.Percentile(costs, 99)
	tp999, _ := stats.Percentile(costs, 99.9)

	var result string
	if tp999/1000 < 1 {
		result = fmt.Sprintf("[%s]: TPS: %.2f, TP99: %.2fus, TP999: %.2fus (b=%d Byte, c=%d, n=%d)",
			title, tps, tp99/1000, tp999/1000, echoSize, concurrent, total)
	} else {
		result = fmt.Sprintf("[%s]: TPS: %.2f, TP99: %.2fms, TP999: %.2fms (b=%d Byte, c=%d, n=%d)",
			title, tps, tp99/1000000, tp999/1000000, echoSize, concurrent, total)
	}
	logInfo(result)
	return nil
}

const blueLayout = "\x1B[1;36;40m%s\x1B[0m"

var infoTitle = "Info: "

func logInfo(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	fmt.Println(infoTitle + s)
}
