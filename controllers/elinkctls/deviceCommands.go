package elinkctls

import (
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type DevCommandController struct {
	ctrl.Controller
}

const pid = "productID"

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
func (this *DevCommandController) Post() {
	pid, err := this.AcquireParamPid()
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pInfo, err := devmodels.LookupProduct(pid)
	if err != nil {
		this.ErrorResponse(200)
		return
	}

	switch pInfo.Types {
	case devmodels.PTypes_Zigbee:
		this.zbDeviceCommandDeal(pid)
	default:
		this.ErrorResponse(303)
	}
}

func (this *DevCommandController) zbDeviceCommandDeal(pid int) {
	var cmdID int

	breq := &ctrl.BaseRequest{}
	bpl := &DevCmdPayload{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &DevCmdRequest{breq, bpl}); err != nil {
		this.ErrorResponse(elink.CodeErrSysInvalidParameter)
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
			this.ErrorResponse(304)
			return
		}
		logs.Debug(cmdID)
		//npis.ZbApps.SendCommand()
	}
}
