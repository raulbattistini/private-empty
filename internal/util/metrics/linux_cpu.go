//go:build linux

package metrics

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/raulbattistini/private-empty/internal/util/logger"
)

type cpuStat struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
}

var (
	lastCPUStat *cpuStat
	lastCPUTime time.Time
)

func readCPUStat() (*cpuStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 8 || fields[0] != "cpu" {
		return nil, os.ErrInvalid
	}

	stat := &cpuStat{}
	stat.user, _ = strconv.ParseUint(fields[1], 10, 64)
	stat.nice, _ = strconv.ParseUint(fields[2], 10, 64)
	stat.system, _ = strconv.ParseUint(fields[3], 10, 64)
	stat.idle, _ = strconv.ParseUint(fields[4], 10, 64)
	stat.iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	stat.irq, _ = strconv.ParseUint(fields[6], 10, 64)
	stat.softirq, _ = strconv.ParseUint(fields[7], 10, 64)
	if len(fields) > 8 {
		stat.steal, _ = strconv.ParseUint(fields[8], 10, 64)
	}

	return stat, nil
}

func getCPUUsagePlatform() float64 {
	curr, err := readCPUStat()
	if err != nil {
		logger.Log().Debug("Failed to read CPU stats: %v", err)
		return 0
	}

	now := time.Now()

	// init
	if lastCPUStat == nil {
		lastCPUStat = curr
		lastCPUTime = now
		return 0
	}

	// 100ms between samples
	duration := now.Sub(lastCPUTime)
	if duration.Milliseconds() < 100 {
		return 0
	}

	// Calculate deltas
	prevIdle := lastCPUStat.idle + lastCPUStat.iowait
	idle := curr.idle + curr.iowait

	prevNonIdle := lastCPUStat.user + lastCPUStat.nice + lastCPUStat.system +
		lastCPUStat.irq + lastCPUStat.softirq + lastCPUStat.steal
	nonIdle := curr.user + curr.nice + curr.system +
		curr.irq + curr.softirq + curr.steal

	prevTotal := prevIdle + prevNonIdle
	total := idle + nonIdle

	totalDelta := total - prevTotal
	idleDelta := idle - prevIdle

	var usage float64
	if totalDelta > 0 {
		usage = float64(totalDelta-idleDelta) / float64(totalDelta) * 100.0
	}

	// next call
	lastCPUStat = curr
	lastCPUTime = now

	return usage
}
