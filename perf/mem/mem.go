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
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultInterval       = time.Second * 3
	defaultUsageThreshold = 1024 // KB
)

type Stats struct {
	Rss int64
	Pss int64
}

type Usage struct {
	MaxRss int64
	AvgRss int64
}

func (u Usage) String() string {
	return fmt.Sprintf("AVG: %d KB, MAX: %d KB", u.AvgRss, u.MaxRss)
}

func RecordUsage(ctx context.Context) (usage Usage, err error) {
	pid := os.Getpid()
	return RecordPidUsage(ctx, pid)
}

func RecordPidUsage(ctx context.Context, pid int) (usage Usage, err error) {
	ticker := time.NewTicker(defaultInterval)
	defer ticker.Stop()
	if pid < 0 {
		pid = os.Getpid()
	}
	var (
		stats  Stats
		bucket []Stats
	)
	for {
		stats, err = getStats(pid)
		if err != nil {
			return
		}
		if stats.Rss > defaultUsageThreshold {
			bucket = append(bucket, stats)
		}

		select {
		case <-ctx.Done():
			return calcStats(bucket), nil
		case <-ticker.C:
		}
	}
}

func calcStats(bucket []Stats) Usage {
	var (
		totalRss int64
		maxRss   int64
	)
	for _, s := range bucket {
		totalRss += s.Rss
		if s.Rss > maxRss {
			maxRss = s.Rss
		}
	}
	return Usage{
		MaxRss: maxRss,
		AvgRss: totalRss / int64(len(bucket)),
	}
}

func parseSmaps(data string) Stats {
	lines := strings.Split(data, "\n")
	reg := regexp.MustCompile("(\\w+): +(\\d+) kB")
	stats := Stats{}
	for _, line := range lines {
		params := reg.FindStringSubmatch(line)
		if len(params) < 3 {
			continue
		}
		field, value := params[1], params[2]
		var err error
		switch strings.ToLower(field) {
		case "rss":
			stats.Rss, err = strconv.ParseInt(value, 10, 64)
		case "pss":
			stats.Pss, err = strconv.ParseInt(value, 10, 64)
		}
		if err != nil {
			panic(fmt.Sprintf("parse mem field: %s failed: %v", field, err))
		}
	}
	return stats
}
