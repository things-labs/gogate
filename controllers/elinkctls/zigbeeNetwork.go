package elinkctls

import (
	"github.com/slzm40/gogate/apps/npis"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/npi"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type ZbNetworkController struct {
	ctrl.Controller
}

// 开启zigbee组网
func (this *ZbNetworkController) Post() {
	var err error
	var ok bool

	if npis.IsNetworkFormation() {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(
			npi.Cms_mode_NetworkSteer) // 开启组网
	} else {
		ok, err = npis.ZbApps.Appcfg_BdbStartCommissioningReq(
			npi.Cms_mode_NetworkFormation | npi.Cms_mode_NetworkSteer) // 建立网络并开启组网
	}
	if err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}
	ctrl.WriteResponse(this.Input, jsoniter.Get(this.Input.Payload, "packetID").ToInt(),
		elink.CodeSuccess, nil)
	npis.SetNetworkSteering(true)
	logs.Debug("elinkctls: zigbee network steering open")
}

// 关闭zigbee组网
func (this *ZbNetworkController) Delete() {
	if ok, err := npis.ZbApps.Zb_PermitJoingReq(0xfffc, 0); err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}
	ctrl.WriteResponse(this.Input, jsoniter.Get(this.Input.Payload, "packetID").ToInt(),
		elink.CodeSuccess, nil)
	npis.SetNetworkSteering(false)
	logs.Debug("elinkctls: zigbee network steering close")
}
