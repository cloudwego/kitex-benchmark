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

/*
#include <unistd.h>
*/
import "C"
import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * man 5 proc
 */
func getPidCPUUsage(ctx context.Context, pid int) chan float64 {
	result := make(chan float64)

	go func() {
		ticker := time.NewTicker(defaultInterval)
		defer func() {
			ticker.Stop()
			close(result)
		}()
		clockTick := float64(C.sysconf(C._SC_CLK_TCK)) // ticks per sec, 机器时钟每秒所走的时钟打点数 // 100.0 //
		lastTime := time.Now()
		lastPIDTime, err := getPIDTime(pid)
		if err != nil {
			log.Fatalf("calculate cpu failed: err=%v", err)
		}
		finished := false
		for !finished {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				finished = true
			}

			nowTime := time.Now()
			pidTime, err := getPIDTime(pid)
			if err != nil {
				log.Fatalf("calculate cpu failed: err=%v", err)
			}

			diffPIDTime := pidTime - lastPIDTime
			period := float64(nowTime.Sub(lastTime).Milliseconds()) / 1000 //seconds in float
			pidCPUUsage := (diffPIDTime) * 100 / (clockTick * period)
			lastPIDTime = pidTime
			lastTime = nowTime

			result <- pidCPUUsage
		}
	}()

	return result
}

// getPIDTime get the sum of utime,stime,cutime,cstime
func getPIDTime(pid int) (float64, error) {
	pidCPUStat := fmt.Sprintf("/proc/%d/stat", pid)
	cpuRet, err := readFile(pidCPUStat)
	if err != nil {
		return 0, err
	}
	items := strings.Split(cpuRet, " ")

	var pidTime float64
	if len(items) < 20 {
		return 0, fmt.Errorf("get cpu info of pid=%d invalid, result=%v", pid, items)
	}
	// 13 is utime:
	//   Amount of time that this process has been scheduled in user mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	//	 This includes guest time, guest_time (time spent running a virtual CPU, see below),
	//   so that applications that are not aware of the guest time field do not lose that time from their calculations.
	// 14 is stime:
	//    Amount of time that this process has been scheduled in kernel mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	// 15 is cutime:
	//    Amount of time that this process's waited-for children have been scheduled in user mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).  (See also times(2).)
	//    This includes guest time,  cguest_time (time spent running a virtual CPU, see below).
	// 16 is cstime:
	//    Amount of time that this process's waited-for children have been scheduled in kernel mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	for i := 13; i <= 16; i++ {
		if t, err := strconv.ParseFloat(items[i], 64); err != nil {
			return 0, err
		} else {
			pidTime += t
		}
	}
	return pidTime, nil
}

func readFile(fileName string) (string, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	b, err := ioutil.ReadAll(fd)
	content := string(b)
	return content, nil
}
