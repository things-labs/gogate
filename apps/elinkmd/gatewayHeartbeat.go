package elinkmd

import (
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/models"
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	Sn string `json:"sn"`
}

// DeviceStatus 设备状态
type DeviceStatus struct {
	CurrentTime   string `json:"currentTime"`
	StartDateTime string `json:"startDateTime"`
	RunningTime   int64  `json:"runningTime"`
	Status        string `json:"status"`
}

// NetInfo 网卡信息
type NetInfo struct {
	MAC string `json:"MAC"`
	Mac string `json:"mac"`
}

// GatewayHeatbeatInfo 网关心跳包信息
type GatewayHeatbeatInfo struct {
	UID          []int64      `json:"uid"`
	DeviceInfo   DeviceInfo   `json:"device_info"`
	DeviceStatus DeviceStatus `json:"device_status"`
	NetInfo      NetInfo      `json:"net_info"`
}

// GetGatewayHeatbeatInfo 心跳包
func GetGatewayHeatbeatInfo(isonline bool) *GatewayHeatbeatInfo {
	status := "online"
	if !isonline {
		status = "offline"
	}

	mac := misc.Mac()
	return &GatewayHeatbeatInfo{
		UID:        models.GetUsers(),
		DeviceInfo: DeviceInfo{Sn: mac},
		DeviceStatus: DeviceStatus{
			CurrentTime:   time.Now().Local().Format("2006-01-02 15:04:05"),
			StartDateTime: misc.SetupTime(),
			RunningTime:   misc.RunningTime(),
			Status:        status,
		},
		NetInfo: NetInfo{
			MAC: misc.MAC(),
			Mac: mac,
		},
	}
}
