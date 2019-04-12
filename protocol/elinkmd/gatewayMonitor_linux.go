// +build linux

package elinkmd

import (
	"fmt"
	"runtime"
	"syscall"
)

type GwMonitor struct {
	SystemMemInfos *syscall.Sysinfo_t
	AppMemInfos    *runtime.MemStats
}

func GatewayMonitors() *GwMonitor {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err != nil {
		fmt.Println("syscall sysinfo failed")
	}
	return &GwMonitor{sysInfo, memStats}
}
