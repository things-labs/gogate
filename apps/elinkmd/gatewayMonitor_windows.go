// +build windows

package elinkmd

import (
	"runtime"
)

// GwMonitor 监控
type GwMonitor struct {
	AppMemInfos *runtime.MemStats
}

// GatewayMonitors 获取监控信息
func GatewayMonitors() (*GwMonitor, error) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	return &GwMonitor{memStats}, nil
}
