package elinkctls

import (
	"github.com/astaxie/beego/logs"
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

type DevCmdPara struct {
	NodeNo  int                    `json:"nodeNo"`
	Command string                 `json:"command"`
	CmdPara map[string]interface{} `json:"cmdPara"`
}

// 设备命令负载
type DevCmdReqPayload struct {
	ProductID int        `json:"productID"`
	Sn        string     `json:"sn"`
	Params    DevCmdPara `json:"params"`
}

// 命令json组合
type DevCmdRequest struct {
	ctrl.BaseRequest
	Payload DevCmdReqPayload `json:"payload,omitempty"`
}

// 下发控制命令
func (this *DevCommandController) Post() {
	pid, err := this.AcquireParamPid()
	if err != nil {
		this.ErrorResponse(elink.CodeErrCommonResourceNotSupport)
		return
	}

	// 确定是否支持此产品
	pInfo, err := devmodels.LookupProduct(pid)
	if err != nil {
		this.ErrorResponse(elink.CodeErrProudctUndefined)
		return
	}

	// 根据产品类型分发命令
	switch pInfo.Types {
	case devmodels.PTypes_Zigbee:
		this.zbDeviceCommandDeal(pid)
	default:
		this.ErrorResponse(elink.CodeErrProudctFeatureUndefined)
	}
}

func (this *DevCommandController) zbDeviceCommandDeal(pid int) {
	code := elink.CodeSuccess
	defer func() { this.ErrorResponse(code) }()

	req := &DevCmdRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}
	logs.Debug("base: %#v", req)
	rpl := req.Payload
	if rpl.Params.NodeNo == ltl.NodeNumReserved {
		dev, err := devmodels.LookupZbDeviceByIeeeAddr(rpl.Sn)
		if err != nil {
			code = elink.CodeErrDeviceNotExist
			return
		}
		switch rpl.Params.Command {
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
			code = elink.CodeErrDeviceCommandNotSupport
			return
		}
		if err != nil {
			code = elink.CodeErrDeviceCommandOperationFailed
			return
		}
		this.WriteResponse(elink.CodeSuccess, nil)
		return
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
