package elinkctls

import (
	"strconv"

	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type DevCommandCtrlController struct {
	ctrl.Controller
}

// 设备命令负载
type DevCmdPayload struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
	NodeNo    int    `json:"nodeNo"`
	Params    struct {
		Command string                 `json:"command"`
		CmdPara map[string]interface{} `json:"cmdPara"`
	} `json:"params"`
}

// 命令json组合
type DevCmdRequest struct {
	*ctrl.BaseRequest
	*DevCmdPayload
}

// 下发控制命令
func (this *DevCommandCtrlController) Post() {
	spid := this.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pid, err := strconv.ParseInt(spid, 10, 0)
	if err != nil { //never happen but deal
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pInfo, exist := devmodels.LookupProduct(int(pid))
	if !exist {
		this.ErrorResponse(200)
		return
	}

	switch pInfo.Types {
	case devmodels.PTypes_Zigbee:
		ZbDeviceCommandDeal(int(pid), this)
	default:
		this.ErrorResponse(303)
	}
}

func ZbDeviceCommandDeal(pid int, dc *DevCommandCtrlController) {
	var cmdID int

	breq := &ctrl.BaseRequest{}
	bpl := &DevCmdPayload{}
	if err := jsoniter.Unmarshal(dc.Input.Payload, &DevCmdRequest{breq, bpl}); err != nil {
		dc.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	//pdtModels.LookupZbDeviceNodeByIN(bpl.Sn, bpl.NodeNo)
	switch pid {
	case devmodels.PID_DZSW01:
		cmd := bpl.Params.Command
		if cmd == "off" {
			cmdID = 0
		} else if cmd == "on" {
			cmdID = 1
		} else if cmd == "toggle" {
			cmdID = 2
		} else {
			dc.ErrorResponse(304)
			return
		}
		logs.Debug(cmdID)
		//npis.ZbApps.SendCommand()
	}
}
