// +build linux

package elinkmd

import (
	"runtime"
	"syscall"
)

// GwMonitor 监控
type GwMonitor struct {
	SystemMemInfos *syscall.Sysinfo_t
	AppMemInfos    *runtime.MemStats
}

// GatewayMonitors 获取监控信息
func GatewayMonitors() (*GwMonitor, error) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err != nil {
		return nil, err
	}
	return &GwMonitor{sysInfo, memStats}, nil
}
