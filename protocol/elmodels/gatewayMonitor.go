package elmodels

import (
	"runtime"
)

func GatewayMonitors() *runtime.MemStats {
	memStats := new(runtime.MemStats)

	runtime.ReadMemStats(memStats)
	return memStats
}
