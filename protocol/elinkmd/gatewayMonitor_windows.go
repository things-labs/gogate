// +build windows

package elinkmd

import (
	"runtime"
)

type GwMonitor struct {
	AppMemInfos *runtime.MemStats
}

func GatewayMonitors() (*GwMonitor, error) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	return &GwMonitor{memStats}, nil
}
