package elinkctls

import (
	"github.com/slzm40/gogate/npis"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type ZbNetworkCtrlController struct {
	ctrl.Controller
}

// 开启zigbee组网
func (this *ZbNetworkCtrlController) Post() {
	var err error
	var ok bool

	if !npis.IsNetworkFormation() {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(0x06) // 建立网络并开启组网
	} else {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(0x04) // 开启组网
	}
	if err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	packetID := jsoniter.Get(this.Input.Payload, "packetID").ToInt()
	ctrl.WriteCtrlResponse(this.Input, packetID, elink.CodeSuccess, nil)
	logs.Debug("elinkctls: zigbee network steering open")
}

// 关闭zigbee组网
func (this *ZbNetworkCtrlController) Delete() {
	if ok, err := npis.ZbApps.Zb_PermitJoingReq(0xfffc, 0); err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	packetID := jsoniter.Get(this.Input.Payload, "packetID").ToInt()
	ctrl.WriteCtrlResponse(this.Input, packetID, elink.CodeSuccess, nil)
	logs.Debug("elinkctls: zigbee network steering close")
}
