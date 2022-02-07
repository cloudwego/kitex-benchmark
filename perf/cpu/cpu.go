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

package cpu

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
)

const (
	defaultInterval       = time.Second * 1
	defaultUsageThreshold = 10 // %
)

type Usage struct {
	Min float64
	Max float64
	Avg float64
	P50 float64
	P90 float64
	P99 float64
}

func (u Usage) String() string {
	return fmt.Sprintf(
		"MIN: %0.2f%%, TP50: %0.2f%%, TP90: %0.2f%%, TP99: %0.2f%%, MAX: %0.2f%%, AVG:%0.2f%%",
		u.Min, u.P50, u.P90, u.P99, u.Max, u.Avg,
	)
}

// RecordUsage return the final Usage when context canceled
func RecordUsage(ctx context.Context) (Usage, error) {
	pid := os.Getpid()
	return RecordUsageWithPid(ctx, pid)
}

// RecordUsageWithPid return the final Usage when context canceled
func RecordUsageWithPid(ctx context.Context, pid int) (usage Usage, err error) {
	if _, err = os.FindProcess(pid); err != nil {
		return
	}

	var cpuUsageList []float64
	var procCpu = sigar.ProcCpu{}
	var ticker = time.NewTicker(defaultInterval)
	defer ticker.Stop()
	for {
		if err = procCpu.Get(pid); err != nil {
			return
		}

		cpuUsage := procCpu.Percent * 100
		if cpuUsage > defaultUsageThreshold {
			cpuUsageList = append(cpuUsageList, cpuUsage)
		}

		select {
		case <-ctx.Done():
			return calcUsage(cpuUsageList), nil
		case <-ticker.C:
		}
	}
}

func calcUsage(stats []float64) Usage {
	if len(stats) == 0 {
		return Usage{}
	}

	sort.Float64s(stats)
	length := len(stats)
	if length > 3 {
		stats = stats[1 : length-1]
		length = length - 2
	}
	fLen := float64(len(stats))
	tp50Index := int(fLen * 0.5)
	tp90Index := int(fLen * 0.9)
	tp99Index := int(fLen * 0.99)

	var usage Usage
	if tp50Index > 0 {
		usage.P50 = stats[tp50Index-1]
	}
	if tp90Index > tp50Index {
		usage.P90 = stats[tp90Index-1]
	} else {
		usage.P90 = usage.P50
	}
	if tp99Index > tp90Index {
		usage.P99 = stats[tp99Index-1]
	} else {
		usage.P99 = usage.P90
	}

	var sum float64
	for _, cost := range stats {
		sum += cost
	}
	usage.Avg = sum / fLen

	usage.Min = stats[0]
	usage.Max = stats[length-1]

	return usage
}
