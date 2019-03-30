package elinkctls

import (
	"errors"

	"github.com/slzm40/gogate/apps/npis"
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/ltl/ltlspec"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

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
	Params    struct {
		NodeNo  int                    `json:"nodeNo"`
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
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeSuccess {
			this.ErrorResponse(code)
		}
	}()
	breq := &ctrl.BaseRequest{}
	bpl := &DevCmdPayload{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &DevCmdRequest{breq, bpl}); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}

	if bpl.Params.NodeNo == ltl.NodeNumReserved {
		dev, err := devmodels.LookupZbDeviceByIeeeAddr(bpl.Sn)
		if err != nil {
			code = elink.CodeErrSysInvalidParameter
			return
		}
		switch bpl.Params.Command {
		case "reset":
			err = npis.ZbApps.SendSpecificCmdBasic(dev.NwkAddr,
				ltlspec.COMMAND_BASIC_REBOOT_DEVICE)
		case "factoryReset":
			err = npis.ZbApps.SendSpecificCmdBasic(dev.NwkAddr,
				ltlspec.COMMAND_BASIC_RESET_FACT_DEFAULT)
		case "identify":
			err = npis.ZbApps.SendSpecificCmdBasic(dev.NwkAddr,
				ltlspec.COMMAND_BASIC_IDENTIFY)
		default:
			err = errors.New("not support")
		}
		if err != nil {
			code = 305
			return
		}
	}

	//	dinfo, err := devmodels.LookupZbDeviceNodeByIN(bpl.Sn, byte(bpl.Params.NodeNo))
	//	if err != nil {
	//		this.ErrorResponse()
	//	}

	//	switch pid {
	//	case devmodels.PID_DZSW01:
	//		cmd := bpl.Params.Command
	//		if cmd == "off" {
	//			cmdID = 0
	//		} else if cmd == "on" {
	//			cmdID = 1
	//		} else if cmd == "toggle" {
	//			cmdID = 2
	//		} else {
	//			this.ErrorResponse(304)
	//			return
	//		}
	//		logs.Debug(cmdID)
	//		//npis.ZbApps.SendCommand()
	//	}
}
