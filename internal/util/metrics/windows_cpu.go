//go:build windows

package metrics

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getSystemTimes = kernel32.NewProc("GetSystemTimes")
	lastIdleTime   uint64
	lastKernelTime uint64
	lastUserTime   uint64
	lastCPUTime    time.Time
)

func fileTimeToUint64(ft *syscall.Filetime) uint64 {
	return uint64(ft.HighDateTime)<<32 + uint64(ft.LowDateTime)
}

func getCPUUsagePlatform() float64 {
	var idleTime, kernelTime, userTime syscall.Filetime

	// idle time
	ret, _, _ := getSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)

	if ret == 0 {
		return 0
	}

	idle := fileTimeToUint64(&idleTime)
	kernel := fileTimeToUint64(&kernelTime)
	user := fileTimeToUint64(&userTime)

	now := time.Now()

	// First call
	if lastIdleTime == 0 {
		lastIdleTime = idle
		lastKernelTime = kernel
		lastUserTime = user
		lastCPUTime = now
		return 0
	}

	duration := now.Sub(lastCPUTime)
	if duration.Milliseconds() < 100 {
		return 0
	}

	idleDelta := idle - lastIdleTime
	kernelDelta := kernel - lastKernelTime
	userDelta := user - lastUserTime

	totalDelta := kernelDelta + userDelta

	var usage float64
	if totalDelta > 0 {
		usage = float64(totalDelta-idleDelta) / float64(totalDelta) * 100.0
	}

	lastIdleTime = idle
	lastKernelTime = kernel
	lastUserTime = user
	lastCPUTime = now

	return usage
}
