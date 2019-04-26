// +build linux

package elinkmd

import (
	"fmt"
	"runtime"
	"syscall"
)

type GwMonitor struct {
	Topic          string `json:"topic,omitempty"`
	SystemMemInfos *syscall.Sysinfo_t
	AppMemInfos    *runtime.MemStats
}

func GatewayMonitors(tp string) *GwMonitor {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err != nil {
		fmt.Println("syscall sysinfo failed")
	}
	return &GwMonitor{tp, sysInfo, memStats}
}
