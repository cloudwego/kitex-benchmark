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
	"sync/atomic"
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
