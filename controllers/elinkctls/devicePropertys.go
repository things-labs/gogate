package elinkctls

import (
	"github.com/slzm40/gomo/ltl"

	"github.com/slzm40/easyjms"
	"github.com/slzm40/gogate/apps/npis"
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gogate/protocol/elmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

	"github.com/json-iterator/go"
)

type DevPropertysController struct {
	ctrl.Controller
}

func (this *DevPropertysController) Get() {
	pid, err := this.AcquireParamPid()
	if err != nil {
		this.ErrorResponse(elink.CodeErrCommonResourceNotSupport)
		return
	}

	pInfo, err := devmodels.LookupProduct(pid)
	if err != nil {
		this.ErrorResponse(elink.CodeErrProudctUndefined)
		return
	}

	switch pInfo.Types {
	case devmodels.PTypes_Zigbee:
		this.zbDevicePropertysDeal(pid)
	default:
		this.ErrorResponse(elink.CodeErrProudctFeatureUndefined)
	}
}

func (this *DevPropertysController) zbDevicePropertysDeal(pid int) {
	req := &elmodels.DevPropRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		this.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	rpl := req.Payload
	jp := easyjms.NewFromMap(rpl.Params)
	if jp.Get("nodeNo").MustInt() == ltl.NodeNumReserved {
		switch jp.Get("types").MustString() {
		case "basic":
			dinfo, err := devmodels.LookupZbDeviceByIeeeAddr(rpl.Sn)
			if err != nil {
				this.ErrorResponse(elink.CodeErrProudctUndefined)
				return
			}

			if err = npis.ZbApps.SendReadReqBasic(dinfo.NwkAddr,
				&elmodels.ItemInfos{
					Pkid:      req.PacketID,
					Client:    this.Input.Client,
					ProductID: rpl.ProductID,
					Sn:        rpl.Sn,
					Tp:        this.Input.Topic,
				}); err != nil {
				this.ErrorResponse(elink.CodeErrDeviceCommandOperationFailed)
				return
			}
		default:
			this.ErrorResponse(elink.CodeErrDevicePropertysNotSupport)
			return
		}

		return
	}
}
