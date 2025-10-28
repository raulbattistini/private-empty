package metrics

// CPU usage percentage (0-100)
func GetCPUUsage() float64 {
	return getCPUUsagePlatform()
}

// memory usage in megabytes
func GetMemoryUsageMB() float64 {
	return getMemoryUsagePlatform()
}

func GetMemoryUsagePercent() float64 {
	return getMemoryUsagePercentPlatform()
}

// GetNetworkUsagePercent returns network usage as percentage
func GetNetworkUsagePercent() float64 {
	return getNetworkUsagePlatform()
}
