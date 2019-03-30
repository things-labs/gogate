package elinkctls

import (
	"errors"

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

type DevProp struct {
	elmodels.BaseSnPayload
	Params map[string]interface{}
}

type DevPropRequest struct {
	*ctrl.BaseRequest
	*DevProp
}

func (this *DevPropertysController) Post() {

}

func (this *DevPropertysController) Get() {
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
		this.zbDevicePropertysDeal(pid)
	default:
		this.ErrorResponse(303)
	}
}

func (this *DevPropertysController) zbDevicePropertysDeal(pid int) error {
	breq := &ctrl.BaseRequest{}
	bpl := &DevProp{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &DevPropRequest{breq, bpl}); err != nil {
		return err
	}

	jp := easyjms.NewFromMap(bpl.Params)
	types := jp.Get("Types").MustString()
	if types == "" {
		return errors.New("error happen")
	}

	switch types {
	case "basic":
		// 忽略节点号
		dinfo, err := devmodels.LookupZbDeviceByIeeeAddr(bpl.Sn)
		if err != nil {
			return err
		}

		err = npis.ZbApps.SendReadReqBasic(dinfo.NwkAddr,
			&elmodels.ItemInfos{Pkid: breq.PacketID})
		if err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}
