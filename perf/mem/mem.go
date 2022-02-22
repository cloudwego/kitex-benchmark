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

package mem

import (
	"context"
	"fmt"
	"os"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
)

const (
	defaultInterval     = time.Second * 3
	defaultRssThreshold = 1024 * 1024 // bytes
)

type Usage struct {
	MaxRss uint64
	AvgRss uint64
}

func (u Usage) String() string {
	return fmt.Sprintf("AVG: %d MB, MAX: %d MB", u.AvgRss/1024/1024, u.MaxRss/1024/1024)
}

// RecordUsage return the final Usage when context canceled
func RecordUsage(ctx context.Context) (usage Usage, err error) {
	pid := os.Getpid()
	return RecordUsageWithPid(ctx, pid)
}

// RecordUsageWithPid return the final Usage when context canceled
func RecordUsageWithPid(ctx context.Context, pid int) (usage Usage, err error) {
	if _, err = os.FindProcess(pid); err != nil {
		return
	}
	var procMem = sigar.ProcMem{}
	var rssList []uint64
	var ticker = time.NewTicker(defaultInterval)
	defer ticker.Stop()
	for {
		if err = procMem.Get(pid); err != nil {
			return
		}
		rss := procMem.Resident
		if rss > defaultRssThreshold {
			rssList = append(rssList, rss)
		}

		select {
		case <-ctx.Done():
			return calcUsage(rssList), nil
		case <-ticker.C:
		}
	}
}

func calcUsage(rssList []uint64) Usage {
	var totalRss, maxRss uint64
	for _, rss := range rssList {
		totalRss += rss
		if rss > maxRss {
			maxRss = rss
		}
	}
	return Usage{
		MaxRss: maxRss,
		AvgRss: totalRss / uint64(len(rssList)),
	}
}
