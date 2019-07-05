package elinkmd

import (
	"runtime"
	"syscall"
)

// GwMonitorInfo 监控信息
type GwMonitorInfo struct {
	SystemMemInfos *syscall.Sysinfo_t
	AppMemInfos    *runtime.MemStats
}

// GetGatewayMonitorInfo 获取监控信息
func GetGatewayMonitorInfo() (*GwMonitorInfo, error) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return nil, err
	}
	return &GwMonitorInfo{sysInfo, memStats}, nil
}
