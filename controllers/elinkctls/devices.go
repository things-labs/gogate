package elinkctls

import (
	"strconv"

	"github.com/slzm40/gogate/models/pdtModels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/slzm40/easyjms"
)

type DevicesCtrlController struct {
	ctrl.CtrlController
}

// 获取设备列表
func (this *DevicesCtrlController) Get() {
	spid := this.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pid, err := strconv.ParseInt(spid, 10, 0)
	if err != nil { //never happen but deal
		this.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pInfo, exist := pdtModels.LookupProduct(int(pid))
	if !exist {
		this.ErrorResponse(200)
		return
	}
	// 根据不同的设备类型分发
	switch pInfo.Types {
	case pdtModels.ProductTypes_General: // 获取通用设备
		getGernalDevices(int(pid), this)
	default:
		this.ErrorResponse(303)
	}
}

// 添加设备
func (this *DevicesCtrlController) Post() {
	dealAddDelGernalDevices(false, this)
}

// 删除设备
func (this *DevicesCtrlController) Delete() {
	dealAddDelGernalDevices(true, this)
}

// 获取通用设备
func getGernalDevices(pid int, dc *DevicesCtrlController) {
	var err error

	code := elink.CodeErrSysInternal
	req := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(dc.Input.Payload, &req); err != nil {
		dc.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	py := []byte{}
	devs := pdtModels.FindGeneralDevice(pid)

	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}

	if py, err = jsoniter.Marshal(struct {
		ProductID int      `json:"productID"`
		Sn        []string `json:"sn"`
	}{pid, sns}); err != nil {
		dc.ErrorResponse(elink.CodeErrSysInternal)
		return
	}
	ctrl.WriteCtrlResponse(dc.Input, req.PacketID, code, py)
}

func dealAddDelGernalDevices(isDel bool, dc *DevicesCtrlController) {
	spid := dc.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		dc.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pid, err := strconv.ParseInt(spid, 10, 0)
	if err != nil { //never happen but deal
		dc.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	pInfo, exist := pdtModels.LookupProduct(int(pid))
	if !exist {
		dc.ErrorResponse(200)
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case pdtModels.ProductTypes_General: // 通用设备处理s
		addDelGernalDevices(isDel, int(pid), dc)
	default:
		dc.ErrorResponse(303)
	}
}

// 添加或删除通用设备
func addDelGernalDevices(isDel bool, pid int, dc *DevicesCtrlController) {
	var sn []string

	code := elink.CodeErrSysInternal
	req := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(dc.Input.Payload, &req); err != nil {
		dc.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	ejs, err := easyjms.NewFromJson(req.Payload)
	if err != nil {
		dc.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	snjs := ejs.Get("sn")
	isArray := snjs.IsArray()
	if isArray {
		sn = snjs.MustStringArray()
	} else {
		sn = append(sn, snjs.MustString())
	}

	snSuc := []string{}
	py := []byte{}
	if len(sn) == 0 {
		dc.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	// 处理要添加或删除的设备
	for _, v := range sn {
		if pdtModels.HasGeneralDevice(pid, v) { // 设备存在
			if isDel {
				if err = pdtModels.DeleteGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		} else { // 设备不存在
			if !isDel {
				if err = pdtModels.CreateGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		}
		snSuc = append(snSuc, v)
	}

	code = elink.CodeSuccess
	if isArray {
		if len(snSuc) == 0 {
			if isDel {
				code = 302
			} else {
				code = 301
			}
		} else if py, err = jsoniter.Marshal(struct {
			ProductID int      `json:"productID"`
			Sn        []string `json:"sn"`
		}{pid, snSuc}); err != nil {
			dc.ErrorResponse(elink.CodeErrSysInternal)
			return
		}
	} else {
		var osn string

		if len(snSuc) > 0 {
			osn = snSuc[0]
		}

		if osn == "" {
			if isDel {
				code = 302
			} else {
				code = 301
			}
		} else if py, err = jsoniter.Marshal(struct {
			ProductID int    `json:"productID"`
			Sn        string `json:"sn"`
		}{pid, osn}); err != nil {
			dc.ErrorResponse(elink.CodeErrSysInternal)
			return
		}
	}

	ctrl.WriteCtrlResponse(dc.Input, req.PacketID, code, py)
}
