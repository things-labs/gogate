package elinkctls

import (
	"github.com/thinkgos/gogate/models/devmodels"
	"github.com/thinkgos/gogate/protocol/elmodels"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/thinkgos/easyjms"
)

type DevicesController struct {
	ctrl.Controller
}

// 获取产品Id下的设备列表
func (this *DevicesController) Get() {
	code := elink.CodeSuccess
	defer func() { this.ErrorResponse(code) }()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrCommonResourceNotSupport
		return
	}

	pInfo, err := devmodels.LookupProduct(pid)
	if err != nil {
		code = elink.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.PTypes_General: // 获取通用设备
		this.getGernalDevices(pid)
	default:
		code = elink.CodeErrProudctFeatureUndefined
	}
}

// 添加设备
func (this *DevicesController) Post() {
	this.dealAddDelGernalDevices(false)
}

// 删除设备
func (this *DevicesController) Delete() {
	this.dealAddDelGernalDevices(true)
}

// 获取通用设备列表
func (this *DevicesController) getGernalDevices(pid int) int {
	devs := devmodels.FindGeneralDevice(pid)
	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}

	py, err := jsoniter.Marshal(elmodels.DevicesInfo{pid, sns})
	if err != nil {
		return elink.CodeErrSysException
	}

	this.WriteResponse(elink.CodeSuccess, py)
	return elink.CodeSuccess
}

func (this *DevicesController) dealAddDelGernalDevices(isDel bool) {
	code := elink.CodeSuccess
	defer func() { this.ErrorResponse(code) }()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrSysException
		return
	}

	pInfo, err := devmodels.LookupProduct(int(pid))
	if err != nil {
		code = elink.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.PTypes_General: // 通用设备处理s
		code = this.addDelGernalDevices(isDel, int(pid))
	default:
		code = elink.CodeErrProudctFeatureUndefined
	}
}

// 添加或删除通用设备
func (this *DevicesController) addDelGernalDevices(isDel bool, pid int) int {
	req := &ctrl.BaseRequest{}
	bpl := &ctrl.BaseRawPayload{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &ctrl.Request{req, bpl}); err != nil {
		return elink.CodeErrSysInvalidParameter

	}
	ejs, err := easyjms.NewFromJson(bpl.Payload)
	if err != nil {
		return elink.CodeErrSysInvalidParameter
	}

	snjs := ejs.Get("sn")
	sn := []string{}
	isArray := snjs.IsArray()
	if isArray {
		sn = snjs.MustStringArray()
	} else if str, err := snjs.String(); err == nil {
		sn = append(sn, str)
	}
	if len(sn) == 0 {
		return elink.CodeErrSysInvalidParameter
	}

	snSuc := []string{}
	py := []byte{}
	// 处理要添加或删除的设备
	for _, v := range sn {
		if devmodels.HasGeneralDevice(pid, v) { // 设备存在
			if isDel {
				if err = devmodels.DeleteGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		} else { // 设备不存在
			if !isDel {
				if err = devmodels.CreateGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		}
		snSuc = append(snSuc, v)
	}

	if isArray {
		if len(snSuc) == 0 {
			return elink.CodeErrDeviceCommandOperationFailed
		}
		if py, err = jsoniter.Marshal(elmodels.DevicesInfo{pid, snSuc}); err != nil {
			return elink.CodeErrSysException
		}
	} else {
		var osn string

		if len(snSuc) > 0 {
			osn = snSuc[0]
		}
		if osn == "" {
			return elink.CodeErrDeviceCommandOperationFailed
		}
		if py, err = jsoniter.Marshal(elmodels.BaseSnPayload{pid, osn}); err != nil {
			return elink.CodeErrSysException
		}
	}

	this.WriteResponse(elink.CodeSuccess, py)
	return elink.CodeSuccess
}
