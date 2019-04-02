package elinkctls

import (
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/npi"
	"github.com/thinkgos/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
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
		this.ErrorResponse(elink.CodeErrSysException)
		return
	}
	npis.SetNetworkSteering(true)
	this.WriteResponse(elink.CodeSuccess, nil)
	logs.Debug("elinkctls: zigbee network steering open")
}

// 关闭zigbee组网
func (this *ZbNetworkController) Delete() {
	if ok, err := npis.ZbApps.Zb_PermitJoingReq(0xfffc, 0); err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysException)
		return
	}
	npis.SetNetworkSteering(false)
	this.WriteResponse(elink.CodeSuccess, nil)
	logs.Debug("elinkctls: zigbee network steering close")
}
