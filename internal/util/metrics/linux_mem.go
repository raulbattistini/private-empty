//go:build linux

package metrics

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/raulbattistini/private-empty/internal/util/logger"
)

func getMemoryUsagePlatform() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.Log().Debug("Failed to read memory stats: %v", err)
		return 0
	}
	defer file.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	if memTotal == 0 {
		return 0
	}

	// memory in MB
	usedKB := memTotal - memAvailable
	usedMB := float64(usedKB) / 1024.0

	return usedMB
}

func getMemoryUsagePercentPlatform() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.Log().Debug("Failed to read memory stats: %v", err)
		return 0
	}
	defer file.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	if memTotal == 0 {
		return 0
	}

	usedKB := memTotal - memAvailable
	percentage := (float64(usedKB) / float64(memTotal)) * 100.0

	return percentage
}
