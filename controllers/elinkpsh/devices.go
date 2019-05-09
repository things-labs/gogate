package elinkpsh

import (
	"github.com/thinkgos/gogate/apps/broad"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/apps/elinkmd"
	"github.com/thinkgos/gomo/elink"
)

// DevSn 设备通知payload
type DevSn struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

// DeviceAnnce 设备加入,离开通知
func DeviceAnnce(pid int, sn string, isjoin bool) error {
	method := elink.MethodDelete
	if isjoin {
		method = elink.MethodPost
	}
	tp := elink.FormatPshTopic(ctrl.ChannelData, elink.FormatResouce(elinkmd.Devices, pid),
		method, elink.MessageTypeAnnce)
	return broad.PublishPyServerJSON(tp, &DevSn{pid, sn})
}
