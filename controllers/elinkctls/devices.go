package elinkctls

import (
	"strconv"

	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/slzm40/easyjms"
)

type DevicesCtrlController struct {
	ctrl.Controller
}

// 获取产品Id下的设备列表
func (this *DevicesCtrlController) Get() {
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeSuccess {
			this.ErrorResponse(code)
		}
	}()

	spid := this.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		code = elink.CodeErrSysInternal
		return
	}

	pid, err := strconv.ParseInt(spid, 10, 0)
	if err != nil { //never happen but deal
		code = elink.CodeErrSysInternal
		return
	}

	pInfo, exist := devmodels.LookupProduct(int(pid))
	if !exist {
		code = 200
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.ProductTypes_General: // 获取通用设备
		getGernalDevices(int(pid), this)
	default:
		code = 202
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

// 获取通用设备列表
func getGernalDevices(pid int, dc *DevicesCtrlController) {
	devs := devmodels.FindGeneralDevice(pid)
	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}

	py, err := jsoniter.Marshal(struct {
		ProductID int      `json:"productID"`
		Sn        []string `json:"sn"`
	}{pid, sns})
	if err != nil {
		dc.ErrorResponse(elink.CodeErrSysInternal)
		return
	}

	packid := jsoniter.Get(dc.Input.Payload, "packetID").ToInt()
	ctrl.WriteCtrlResponse(dc.Input, packid, elink.CodeSuccess, py)
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

	pInfo, exist := devmodels.LookupProduct(int(pid))
	if !exist {
		dc.ErrorResponse(200)
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case devmodels.ProductTypes_General: // 通用设备处理s
		addDelGernalDevices(isDel, int(pid), dc)
	default:
		dc.ErrorResponse(202)
	}
}

// 添加或删除通用设备
func addDelGernalDevices(isDel bool, pid int, dc *DevicesCtrlController) {
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeSuccess {
			dc.ErrorResponse(code)
		}
	}()

	req := &ctrl.BaseRequest{}
	bpl := &ctrl.BasePayload{}
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
		if py, err = jsoniter.Marshal(struct {
			ProductID int      `json:"productID"`
			Sn        []string `json:"sn"`
		}{pid, snSuc}); err != nil {
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
		if py, err = jsoniter.Marshal(struct {
			ProductID int    `json:"productID"`
			Sn        string `json:"sn"`
		}{pid, osn}); err != nil {
			code = elink.CodeErrSysInternal
			return
		}
	}

	ctrl.WriteCtrlResponse(dc.Input, req.PacketID, code, py)
}
