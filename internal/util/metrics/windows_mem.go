//go:build windows

package metrics

import (
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var (
	kernelWin32          = syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx = kernelWin32.NewProc("GlobalMemoryStatusEx")
)

func getMemoryUsagePlatform() float64 {
	var memStatus memoryStatusEx
	memStatus.dwLength = uint32(unsafe.Sizeof(memStatus))

	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return 0
	}

	usedBytes := memStatus.ullTotalPhys - memStatus.ullAvailPhys
	usedMB := float64(usedBytes) / 1024.0 / 1024.0

	return usedMB
}

func getMemoryUsagePercentPlatform() float64 {
	var memStatus memoryStatusEx
	memStatus.dwLength = uint32(unsafe.Sizeof(memStatus))

	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return 0
	}

	// dwMemoryLoad gives a simple result here
	return float64(memStatus.dwMemoryLoad)
}
