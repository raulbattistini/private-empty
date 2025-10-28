//go:build windows

package metrics

import (
	"math/rand"
)

//  to properly calculate windows network usage it requires complex WMI queries
// simulated data to sample
// possible impl: GetIfTable2 from iphlpapi.dll

var (
	baselineNet = 1024 * 1024 * (1 + rand.Float64()*5)
	lastNet     = baselineNet
)

func getNetworkUsagePlatform() float64 {
	if rand.Float64() < 0.1 {
		lastNet = baselineNet * (2 + rand.Float64()*3)
	} else {
		change := (rand.Float64() - 0.5) * baselineNet * 0.3
		lastNet += change
		lastNet = lastNet*0.8 + baselineNet*0.2
	}

	if lastNet < 0 {
		lastNet = 0
	}
	if lastNet > baselineNet*10 {
		lastNet = baselineNet * 10
	}

	const gigabitInBytes = 1000 * 1000 * 1000 / 8
	return (lastNet / gigabitInBytes) * 100
}
