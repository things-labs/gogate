package elinkpsh

import (
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
)

type DevSnPy struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

// 设备加入或离开通知
func DeviceAnnce(pid int, sn string, isjoin bool) error {
	v, err := jsoniter.Marshal(DevSnPy{pid, sn})
	if err != nil {
		return err
	}
	method := elink.MethodDelete
	if isjoin {
		method = elink.MethodPost
	}

	return ctrl.Publish(elink.FormatResouce(elinkmd.Devices, pid),
		method, elink.MessageTypeAnnce, v)
}
