// +build windows

package elinkmd

import (
	"runtime"
)

type GwMonitor struct {
	Topic       string `json:"topic,omitempty"`
	AppMemInfos *runtime.MemStats
}

func GatewayMonitors(tp string) *GwMonitor {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	return &GwMonitor{tp, memStats}
}
