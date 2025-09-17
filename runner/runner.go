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
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/juju/ratelimit"
)

// 为了流量更均匀, 时间间隔设置为 10ms
const window = 10 * time.Millisecond

// 单次测试
type RunOnce func() error

type Runner struct {
	counters []*Counter // 计数器
	timer    Timer      // 计时器
}

func NewRunner() *Runner {
	r := &Runner{}
	if qps == 0 {
		r.timer = NewTimer(time.Microsecond)
	} else {
		r.timer = NewTimer(0)
	}
	return r
}

func (r *Runner) benching(onceFn RunOnce, concurrent, qps int, duration int64) {
	var wg sync.WaitGroup
	wg.Add(concurrent)
	r.counters = nil
	var qpsLimiter *ratelimit.Bucket
	if qps > 0 {
		qpsLimiter = ratelimit.NewBucketWithRate(float64(qps), int64(concurrent))
	}
	var stopped uint32
	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		atomic.StoreUint32(&stopped, 1)
	})
	for i := 0; i < concurrent; i++ {
		c := NewCounter()
		r.counters = append(r.counters, c)
		go func(c *Counter) {
			defer wg.Done()
			for {
				if atomic.LoadUint32(&stopped) == 1 {
					break
				}
				if qpsLimiter != nil {
					qpsLimiter.Wait(1)
				}
				begin := r.timer.Now()
				err := onceFn()
				end := r.timer.Now()
				if err != nil {
					if errors.Is(err, kerrors.ErrCircuitBreak) {
						klog.Warnf("No.%d request failed: %v, circuit break happens, stop test!", c.Total, err)
						break
					}
					klog.Warnf("No.%d request failed: %v", c.Total, err)
				}
				cost := end - begin
				c.AddRecord(err, cost)
			}
		}(c)
	}
	wg.Wait()
}

func (r *Runner) Warmup(onceFn RunOnce, concurrent, qps int, duration int64) {
	r.benching(onceFn, concurrent, qps, duration)
}

// 并发测试
func (r *Runner) Run(title string, onceFn RunOnce, concurrent, qps int, duration int64, echoSize, sleepTime int) {
	logInfo(
		"%s start benching [%s], concurrent: %d, qps: %d, duration: %ds, sleep: %d",
		"["+title+"]", time.Now().String(), concurrent, qps, duration, sleepTime,
	)

	start := r.timer.Now()
	r.benching(onceFn, concurrent, qps, duration)
	stop := r.timer.Now()
	var totalCounter Counter
	for _, c := range r.counters {
		totalCounter.Total += c.Total
		totalCounter.Failed += c.Failed
		totalCounter.costs = append(totalCounter.costs, c.costs...)
	}
	totalCounter.Report(title, stop-start, concurrent, duration, echoSize)
}
