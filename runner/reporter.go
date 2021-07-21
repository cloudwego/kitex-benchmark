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
	"time"

	"github.com/montanaflynn/stats"
)

func (c *Counter) Report(title string, totalns int64, concurrent int, total int64, echoSize int) error {
	ms, sec := int64(time.Millisecond), int64(time.Second)
	logInfo("took %d ms for %d requests", totalns/ms, c.Total)
	logInfo("requests total: %d, failed: %d", c.Total, c.Failed)

	var tps float64
	if totalns < sec {
		tps = float64(c.Total*sec) / float64(totalns)
	} else {
		tps = float64(c.Total) / (float64(totalns) / float64(sec))
	}

	var costs = make([]float64, len(c.costs))
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
	logInfo(blueString(result))
	return nil
}

const blue_layout = "\x1B[1;36;40m%s\x1B[0m"

var infoTitle = blueString("Info: ")

func blueString(s string) string {
	return fmt.Sprintf(blue_layout, s)
}

func logInfo(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	fmt.Println(infoTitle + s)
}
