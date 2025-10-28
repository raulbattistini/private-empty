package util

import (
	"runtime"
	"time"
)

type cpuSample struct {
	idle  uint64
	total uint64
}

var lastSample cpuSample
var lastSampleTime time.Time

func init() {
	lastSample = takeCPUSample()
	lastSampleTime = time.Now()
}

func takeCPUSample() cpuSample {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	numCPU := runtime.NumCPU()
	numGoroutines := runtime.NumGoroutine()

	return cpuSample{
		idle:  uint64(numCPU * 100),
		total: uint64(numGoroutines * 10),
	}
}

func GetCPUUsage() float64 {
	now := time.Now()
	current := takeCPUSample()

	deltaTime := now.Sub(lastSampleTime).Seconds()
	if deltaTime < 0.1 {
		return 0
	}

	deltaTotal := float64(current.total - lastSample.total)
	deltaIdle := float64(current.idle - lastSample.idle)

	usage := (deltaTotal / (deltaTotal + deltaIdle)) * 100

	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}

	lastSample = current
	lastSampleTime = now

	return usage
}
