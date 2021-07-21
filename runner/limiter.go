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
	"sync"
	"sync/atomic"
	"time"
)

// QPS 限制器
type Limiter struct {
	sync.Once
	qps    int64 // 总 qps
	now    int64
	window time.Duration
	limit  int64 // qps max per window
}

// 公用全局 Timer reset 减少开销, 但 AddNotify 无法回收, 因此不能频繁创建 Limiter
// must use *Limiter, 否则配置会被复制导致失效
func NewLimiter(maxQps int64, window time.Duration) *Limiter {
	limiter := &Limiter{
		window: window,
	}
	limiter.Reset(maxQps)
	limiter.refresh()
	return limiter
}

// has race here but will not appear
func (l *Limiter) Reset(maxQps int64) {
	l.limit = maxQps / int64(time.Second/window)
	l.reset()
}

// true 代表 超过上限, maxQps <= 0表示没有上限
func (l *Limiter) QpsOverrun() bool {
	if l.limit <= 0 {
		return false
	}
	if cur := atomic.AddInt64(&l.qps, 1); cur > l.limit {
		return true
	}
	return false
}

func (l *Limiter) reset() {
	atomic.StoreInt64(&l.qps, 0)
}

func (l *Limiter) refresh() {
	l.Do(func() {
		go func() {
			for range time.Tick(l.window) {
				l.reset()
			}
		}()
	})
}
