package elinkctls

import (
	"github.com/slzm40/gogate/npis"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type ZbNetworkCtrlController struct {
	ctrl.CtrlController
}

// 开启zigbee网络
func (this *ZbNetworkCtrlController) Post() {
	var err error
	var ok bool

	packetID := jsoniter.Get(this.Input.Payload, "packetID").ToInt()

	if !npis.ZbApps.IsNetworkFormation {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(0x06) // 建立网络并开启组网
	} else {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(0x04) // 开启组网
	}

	if err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	code := elink.CodeSuccess
	ctrl.WriteCtrlResponse(this.Input, packetID, code, nil)
	logs.Debug("elinkctls: zigbee network steering open")
}

// 关闭zigbee网络
func (this *ZbNetworkCtrlController) Delete() {
	packetID := jsoniter.Get(this.Input.Payload, "packetID").ToInt()

	ok, err := npis.ZbApps.Zb_PermitJoingReq(0xfffc, 0)
	if err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	code := elink.CodeSuccess
	ctrl.WriteCtrlResponse(this.Input, packetID, code, nil)
	logs.Debug("elinkctls: zigbee network steering close")
}
