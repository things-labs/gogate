package elinkmd

import (
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/utils"
)

type DeviceInfo struct {
	Sn string `json:"sn"`
}
type DeviceStatus struct {
	CurrentTime   string `json:"currentTime"`
	StartDateTime string `json:"startDateTime"`
	RunningTime   string `json:"runningTime"`
	Status        string `json:"status"`
}
type NetInfo struct {
	MAC string `json:"MAC"`
	Mac string `json:"mac"`
}

type GatewayHeatbeat struct {
	Uid          []int64      `json:"uid"`
	DeviceInfo   DeviceInfo   `json:"device_info"`
	DeviceStatus DeviceStatus `json:"device_status"`
	NetInfo      NetInfo      `json:"net_info"`
}

func GatewayHeatbeats(isonline bool) *GatewayHeatbeat {
	status := "online"
	if !isonline {
		status = "offline"
	}
	mac := misc.Mac()
	return &GatewayHeatbeat{
		Uid:        models.GetUsers(),
		DeviceInfo: DeviceInfo{Sn: mac},
		DeviceStatus: DeviceStatus{
			CurrentTime:   time.Now().Local().Format("2006-01-02 15:04:05"),
			StartDateTime: utils.SetupTime(),
			RunningTime:   utils.RunningTime(),
			Status:        status,
		},
		NetInfo: NetInfo{
			MAC: misc.MAC(),
			Mac: mac,
		},
	}
}
