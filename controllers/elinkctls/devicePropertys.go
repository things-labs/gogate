package elinkctls

import (
	"github.com/thinkgos/easyjms"
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/middle/ewait"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/protocol/limp"

	"github.com/json-iterator/go"
)

type DevPropReqPy struct {
	ProductID int                    `json:"productID"`
	Sn        string                 `json:"sn"`
	NodeNo    int                    `json:"nodeNo,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

type DevPropRequest struct {
	ctrl.BaseRequest
	Payload DevPropReqPy `json:"payload,omitempty"`
}

type DevPropRspPy struct {
	ProductID int         `json:"productID"`
	Sn        string      `json:"sn"`
	NodeNo    int         `json:"nodeNo,omitempty"`
	Data      interface{} `json:"data"`
}

type DevPropertysController struct {
	ctrl.Controller
}

func (this *DevPropertysController) Get() {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrCommonResourceNotSupport
		return
	}
	// 确定产品Id是否注册过
	pInfo, err := models.LookupProduct(pid)
	if err != nil {
		code = elink.CodeErrProudctUndefined
		return
	}

	switch pInfo.Types {
	case models.PTypes_Zigbee:
		code = this.zbDevicePropertysGet(pid)
	default:
		code = elink.CodeErrProudctFeatureUndefined
	}
}

func (this *DevPropertysController) zbDevicePropertysGet(pid int) int {
	req := &DevPropRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		return elink.CodeErrSysInvalidParameter
	}

	rpl := req.Payload
	jp := easyjms.NewFromMap(rpl.Params)
	if rpl.NodeNo == ltl.NodeNumReserved {
		switch jp.Get("types").MustString() {
		case "basic":
			dinfo, err := models.LookupZbDeviceByIeeeAddr(rpl.Sn)
			if err != nil {
				return elink.CodeErrDeviceNotExist
			}
			id := ewait.ObainID()
			if err = npis.ZbApps.SendReadReqBasic(dinfo.NwkAddr, id); err != nil {
				return elink.CodeErrDeviceCommandOperationFailed
			}
			v, ok := ewait.Add(id).Wait()
			if !ok {
				return elink.CodeErrDeviceCommandOperationFailed
			}
			item, ok := v.(*limp.GenerlBasicAttribute)
			if !ok {
				return elink.CodeErrDeviceCommandOperationFailed
			}
			out, err := jsoniter.Marshal(DevPropRspPy{rpl.ProductID, rpl.Sn, rpl.NodeNo, item})
			if err != nil {
				return elink.CodeErrSysException
			}
			this.WriteResponsePy(elink.CodeSuccess, out)
		default:
			return elink.CodeErrDevicePropertysNotSupport
		}
		return elink.CodeSuccess
	}

	return elink.CodeSuccess
}
