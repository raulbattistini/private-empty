package util

import (
	"math/rand"
)

var (
	baselineNet = 1024 * 1024 * (1 + rand.Float64()*5) // conservative 1-6 MB/s baseline
	lastNet     = baselineNet
)

func GetNetworkUsage() float64 {
	// simulating burts
	if rand.Float64() < 0.1 {
		lastNet = baselineNet * (2 + rand.Float64()*3) // if spikes
	} else {
		// normal variation
		change := (rand.Float64() - 0.5) * baselineNet * 0.3
		lastNet += change

		// a mean reversion tweak
		lastNet = lastNet*0.8 + baselineNet*0.2
	}

	// making it realistical
	if lastNet < 0 {
		lastNet = 0
	}
	if lastNet > baselineNet*10 {
		lastNet = baselineNet * 10
	}

	return lastNet
}

// network usage as percentage of a (common) 100 Mbps link
func GetNetworkUsagePercent() float64 {
	const maxBandwidth = 100 * 1024 * 1024 / 8 // 100 Mbps translated into bytes/sec
	return (GetNetworkUsage() / maxBandwidth) * 100
}
