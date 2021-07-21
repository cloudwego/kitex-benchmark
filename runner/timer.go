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

func NewTimer(window time.Duration) *Timer {
	t := &Timer{window: window}
	t.refresh()
	return t
}

// 全局 Timer, 共享时间周期, 并在到期时执行回调
type Timer struct {
	sync.Once
	now    int64
	window time.Duration
	notify []func(now time.Time)
}

// refresh time
func (t *Timer) refresh() {
	t.Do(func() {
		atomic.StoreInt64(&t.now, time.Now().UnixNano())
		go func() {
			for now := range time.Tick(t.window) {
				atomic.StoreInt64(&t.now, now.UnixNano())
			}
		}()
	})
}

func (t *Timer) Window() time.Duration {
	return t.window
}

// Timer 为共享计时器, 减少系统时间调用
func (t *Timer) Now() int64 {
	return atomic.LoadInt64(&t.now)
}
