package elinkpsh

import (
	"github.com/thinkgos/gogate/apps/broad"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"
)

type DevSnPy struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

// 设备加入或离开通知
func DeviceAnnce(pid int, sn string, isjoin bool) error {
	method := elink.MethodDelete
	if isjoin {
		method = elink.MethodPost
	}
	tp := elink.FormatPshTopic(ctrl.ChannelData, elink.FormatResouce(elinkmd.Devices, pid),
		method, elink.MessageTypeAnnce)
	return broad.PublishPyServerJSON(tp, DevSnPy{pid, sn})
}
