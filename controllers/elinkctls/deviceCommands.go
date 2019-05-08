package elinkctls

import (
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/ltl/ltlspec"

	jsoniter "github.com/json-iterator/go"
)

type DevCmdPara struct {
	Command string                 `json:"command"`
	CmdPara map[string]interface{} `json:"cmdPara"`
}

// 设备命令负载
type DevCmdReqPy struct {
	ProductID int        `json:"productID"`
	Sn        string     `json:"sn"`
	NodeNo    int        `json:"nodeNo"`
	Params    DevCmdPara `json:"params"`
}

// 命令json组合
type DevCmdRequest struct {
	ctrl.BaseRequest
	Payload DevCmdReqPy `json:"payload,omitempty"`
}

type DevCommandController struct {
	ctrl.Controller
}

// 下发控制命令
func (this *DevCommandController) Post() {
	pid, err := this.AcquireParamPid()
	if err != nil {
		this.ErrorResponse(elink.CodeErrCommonResourceNotSupport)
		return
	}

	// 确定是否支持此产品
	pInfo, err := models.LookupProduct(pid)
	if err != nil {
		this.ErrorResponse(elink.CodeErrProudctUndefined)
		return
	}

	// 根据产品类型分发命令
	switch pInfo.Types {
	case models.PTypes_Zigbee:
		this.zbDeviceCommandDeal(pid)
	default:
		this.ErrorResponse(elink.CodeErrProudctFeatureUndefined)
	}
}

func (this *DevCommandController) zbDeviceCommandDeal(pid int) {
	var cmdID byte

	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	req := &DevCmdRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}
	rpl := req.Payload
	if rpl.NodeNo == ltl.NodeNumReserved {
		dev, err := models.LookupZbDeviceByIeeeAddr(rpl.Sn)
		if err != nil {
			code = elink.CodeErrDeviceNotExist
			return
		}
		switch rpl.Params.Command {
		case "reboot":
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
		this.WriteResponsePyServerJSON(elink.CodeSuccess, nil)
		return
	}

	dinfo, err := models.LookupZbDeviceNodeByIN(rpl.Sn, byte(rpl.NodeNo))
	if err != nil {
		code = elink.CodeErrDeviceNotExist
		return
	}

	switch pid {
	case models.PID_DZSW01, models.PID_DZSW02, models.PID_DZSW03:
		cmd := rpl.Params.Command
		if cmd == "off" {
			cmdID = 0
		} else if cmd == "on" {
			cmdID = 1
		} else if cmd == "toggle" {
			cmdID = 2
		} else {
			code = elink.CodeErrDeviceCommandNotSupport
			return
		}
		err := npis.ZbApps.SendSpecificCmd(dinfo.GetNwkAddr(), ltl.TrunkID_GeneralOnoff,
			byte(rpl.NodeNo), ltl.LTL_FRAMECTL_CLIENT_SERVER_DIR, ltl.RESPONSETYPE_NO, cmdID, nil, nil)
		if err != nil {
			code = elink.CodeErrDeviceCommandOperationFailed
			return
		}
		this.WriteResponsePyServerJSON(elink.CodeSuccess, nil)
	default:
		code = elink.CodeErrProudctFeatureUndefined
		return
	}
}
