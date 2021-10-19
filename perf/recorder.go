//
// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package perf

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/cloudwego/kitex-benchmark/perf/cpu"
	"github.com/cloudwego/kitex-benchmark/perf/mem"
)

type Recorder struct {
	name     string
	finish   func()
	waiter   sync.WaitGroup
	CpuUsage cpu.Usage
	MemUsage mem.Usage
}

func NewRecorder(name string) *Recorder {
	return &Recorder{
		name: name,
	}
}

func (r *Recorder) Begin() {
	r.Reset()

	ctx, finish := context.WithCancel(context.Background())
	r.finish = finish
	r.waiter.Add(2)
	go func() {
		defer r.waiter.Done()
		var err error
		r.CpuUsage, err = cpu.RecordUsage(ctx)
		if err != nil {
			log.Fatalf("recording cpu usage failed: %v", err)
		}
	}()
	go func() {
		defer r.waiter.Done()
		var err error
		r.MemUsage, err = mem.RecordUsage(ctx)
		if err != nil {
			log.Fatalf("recording mem usage failed: %v", err)
		}
	}()
}

func (r *Recorder) End() {
	r.finish()
	r.waiter.Wait()
}

func (r *Recorder) Reset() {
	r.finish = nil
	r.waiter = sync.WaitGroup{}
	r.CpuUsage = cpu.Usage{}
	r.MemUsage = mem.Usage{}
}

func (r *Recorder) ReportString() string {
	output := ""
	output += fmt.Sprintf("[%s] CPU Usage: %s\n", r.name, r.CpuUsage)
	output += fmt.Sprintf("[%s] Mem Usage: %s\n", r.name, r.MemUsage)
	return output
}

func (r *Recorder) Report() {
	fmt.Print(r.ReportString())
}
