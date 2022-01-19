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

//go:build darwin || netbsd || freebsd || openbsd || dragonfly
// +build darwin netbsd freebsd openbsd dragonfly

package cpu

import (
	"context"
	"math/rand"
	"time"
)

func getPidCPUUsage(ctx context.Context, pid int) chan float64 {
	result := make(chan float64, 1)
	go func() {
		ticker := time.NewTicker(defaultInterval)
		defer func() {
			ticker.Stop()
			close(result)
		}()
		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}

			result <- rand.Float64() * 100
		}
	}()

	return result
}
