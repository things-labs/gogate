package elinkctls

import (
	"github.com/thinkgos/elink"

	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/models"
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
		code = elink.CodeErrSysResourceNotSupport
		return
	}
	// 确定产品Id是否注册过
	pInfo, err := models.LookupProduct(pid)
	if err != nil {
		code = ctrl.CodeErrProudctUndefined
		return
	}

	switch pInfo.Types {
	case models.PTypesZigbee:
		code = this.zbDevicePropertysGet(pid)
	default:
		code = ctrl.CodeErrProudctFeatureUndefined
	}
}

func (this *DevPropertysController) zbDevicePropertysGet(pid int) int {
	// req := &DevPropRequest{}
	// if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
	// 	return elink.CodeErrSysInvalidParameter
	// }
	//
	// rpl := req.Payload

	// jp := easyjms.NewFromMap(rpl.Params)
	// if rpl.NodeNo == ltl.NodeNumReserved {
	// 	switch jp.Get("types").MustString() {
	// 	case "basic":
	// 		dinfo, err := models.LookupZbDeviceByIeeeAddr(rpl.Sn)
	// 		if err != nil {
	// 			return ctrl.CodeErrDeviceNotExist
	// 		}
	// 		id := this.SyncManage.ObainID()
	// 		if err = npis.ZbApps.SendReadReqBasic(dinfo.NwkAddr, id); err != nil {
	// 			return ctrl.CodeErrDeviceCommandOperationFailed
	// 		}
	//
	// 		v, ok := this.SyncManage.Wait(id)
	// 		if !ok {
	// 			return ctrl.CodeErrDeviceCommandOperationFailed
	// 		}
	// 		item, ok := v.(*limp.GenerlBasicAttribute)
	// 		if !ok {
	// 			return ctrl.CodeErrDeviceCommandOperationFailed
	// 		}
	//
	// 		err = this.WriteResponsePyServerJSON(elink.CodeSuccess,
	// 			&DevPropRspPy{rpl.ProductID, rpl.Sn, rpl.NodeNo, item})
	// 		if err != nil {
	// 			return elink.CodeErrSysException
	// 		}
	// 	default:
	// 		return ctrl.CodeErrDevicePropertysNotSupport
	// 	}
	// 	return elink.CodeSuccess
	// }

	return elink.CodeSuccess
}
