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

type netStat struct {
	rxBytes uint64
	txBytes uint64
}

var (
	lastNetStat *netStat
	lastNetTime time.Time
)

func readNetStat() (*netStat, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	total := &netStat{}
	scanner := bufio.NewScanner(file)

	// Skip first 2 header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 10 {
			continue
		}

		// Skip loopback
		if strings.Contains(fields[0], "lo:") {
			continue
		}

		// RxBytes is field 1, TxBytes is field 9
		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)

		total.rxBytes += rx
		total.txBytes += tx
	}

	return total, scanner.Err()
}

func getNetworkUsagePlatform() float64 {
	curr, err := readNetStat()
	if err != nil {
		logger.Log().Debug("Failed to read network stats: %v", err)
		return 0
	}

	now := time.Now()

	if lastNetStat == nil {
		lastNetStat = curr
		lastNetTime = now
		return 0
	}

	duration := now.Sub(lastNetTime).Seconds()
	if duration < 0.1 {
		return 0
	}

	rxDelta := curr.rxBytes - lastNetStat.rxBytes
	txDelta := curr.txBytes - lastNetStat.txBytes
	totalDelta := rxDelta + txDelta

	bytesPerSec := float64(totalDelta) / duration

	lastNetStat = curr
	lastNetTime = now

	// Return as percentage of 1 Gbps
	const gigabitInBytes = 1000 * 1000 * 1000 / 8
	return (bytesPerSec / gigabitInBytes) * 100
}
