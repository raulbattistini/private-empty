package util

import "runtime"

var RunningOs = runtime.GOOS

type OSName string

const (
	UNKNOWN OSName = "OSName"
	LINUX   OSName = "linux"
	DARWIN  OSName = "darwin"
	WINDOWS OSName = "windows"
)

func ParseRuntimeOS() OSName {
	osStr := runtime.GOOS

	switch osStr {
	case "linux":
		return LINUX
	case "darwin":
		return DARWIN
	case "windows":
		return WINDOWS
	default:
		return UNKNOWN
	}
}
