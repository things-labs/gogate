package elinkctls

import (
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gogate/protocol/elmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/slzm40/easyjms"
)

type DevicesController struct {
	ctrl.Controller
}

// 获取产品Id下的设备列表
func (this *DevicesController) Get() {
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeSuccess {
			this.ErrorResponse(code)
		}
	}()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrSysInternal
		return
	}

	pInfo, err := devmodels.LookupProduct(pid)
	if err != nil {
		code = 200
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.PTypes_General: // 获取通用设备
		getGernalDevices(pid, this)
	default:
		code = 202
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
func getGernalDevices(pid int, dc *DevicesController) {
	devs := devmodels.FindGeneralDevice(pid)
	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}

	py, err := jsoniter.Marshal(elmodels.DevicesInfo{pid, sns})
	if err != nil {
		dc.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	packid := jsoniter.Get(dc.Input.Payload, "packetID").ToInt()
	ctrl.WriteResponse(dc.Input, packid, elink.CodeSuccess, py)
}

func (this *DevicesController) dealAddDelGernalDevices(isDel bool) {
	pid, err := this.AcquireParamPid()
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pInfo, err := devmodels.LookupProduct(int(pid))
	if err != nil {
		this.ErrorResponse(200)
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.PTypes_General: // 通用设备处理s
		addDelGernalDevices(isDel, int(pid), this)
	default:
		this.ErrorResponse(202)
	}
}

// 添加或删除通用设备
func addDelGernalDevices(isDel bool, pid int, dc *DevicesController) {
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeSuccess {
			dc.ErrorResponse(code)
		}
	}()

	req := &ctrl.BaseRequest{}
	bpl := &ctrl.BaseRawPayload{}
	if err := jsoniter.Unmarshal(dc.Input.Payload, &ctrl.Request{req, bpl}); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}
	ejs, err := easyjms.NewFromJson(bpl.Payload)
	if err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
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
		code = elink.CodeErrSysInvalidParameter
		return
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
			code = 301
			return
		}
		if py, err = jsoniter.Marshal(elmodels.DevicesInfo{pid, snSuc}); err != nil {
			code = elink.CodeErrSysInternal
			return
		}
	} else {
		var osn string

		if len(snSuc) > 0 {
			osn = snSuc[0]
		}
		if osn == "" {
			code = 301
			return
		}
		if py, err = jsoniter.Marshal(elmodels.BaseSnPayload{pid, osn}); err != nil {
			code = elink.CodeErrSysInternal
			return
		}
	}

	ctrl.WriteResponse(dc.Input, req.PacketID, code, py)
}
