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

import "sync/atomic"

// Pool implements a round-robin balanced object pool
type Pool struct {
	clients []interface{}
	size    uint64
	cursor  uint64
}

func NewPool(factory func() interface{}, size int) *Pool {
	p := &Pool{
		clients: make([]interface{}, 0, size),
		size:    uint64(size),
	}
	for i := 0; i < size; i++ {
		p.clients = append(p.clients, factory())
	}
	return p
}

func (p *Pool) Get() interface{} {
	return p.clients[atomic.AddUint64(&p.cursor, 1)%p.size]
}
