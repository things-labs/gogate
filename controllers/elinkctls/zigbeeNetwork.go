package elinkctls

import (
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/npi"

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
	logs.Debug("elinkctls: zigbee network steering open")

	err = this.WriteResponse(elink.CodeSuccess, nil)
	if err != nil {
		logs.Error(err)
	}

}

// 关闭zigbee组网
func (this *ZbNetworkController) Delete() {
	ok, err := npis.ZbApps.Zb_PermitJoingReq(0xfffc, 0)
	if err != nil && !ok {
		this.ErrorResponse(elink.CodeErrSysException)
		return
	}
	npis.SetNetworkSteering(false)
	logs.Debug("elinkctls: zigbee network steering close")

	err = this.WriteResponse(elink.CodeSuccess, nil)
	if err != nil {
		logs.Error(err)
	}

}
