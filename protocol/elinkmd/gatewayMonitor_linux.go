// +build linux

package elinkmd

import (
	"fmt"
	"runtime"
	"syscall"
)

type GwMonitor struct {
	AppMemInfos    *runtime.MemStats
	SystemMemInfos *syscall.Sysinfo_t
}

func GatewayMonitors() *GwMonitor {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err != nil {
		fmt.Println("syscall sysinfo failed")
	}
	return &GwMonitor{memStats, sysInfo}
}
