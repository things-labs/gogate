package elinkctls

import (
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

type DevSnPy struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

type DevMultiSnPy struct {
	ProductID int      `json:"productID"`
	Sn        []string `json:"sn"`
}

type DevMultiSnRequest struct {
	ctrl.BaseRequest
	Payload DevMultiSnPy `json:"payload,omitempty"`
}

type DevicesController struct {
	ctrl.Controller
}

// 获取产品Id下的设备列表
func (this *DevicesController) Get() {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrCommonResourceNotSupport
		return
	}

	pInfo, err := models.LookupProduct(pid)
	if err != nil {
		code = elink.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case models.PTypes_General: // 获取通用设备
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
	devs := models.FindGeneralDevice(pid)
	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}

	this.WriteResponsePyServerJSON(elink.CodeSuccess, &DevMultiSnPy{pid, sns})
	return elink.CodeSuccess
}

func (this *DevicesController) dealAddDelGernalDevices(isDel bool) {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrSysException
		return
	}

	pInfo, err := models.LookupProduct(int(pid))
	if err != nil {
		code = elink.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case models.PTypes_General: // 通用设备处理s
		code = this.addDelGernalDevices(isDel, int(pid))
	default:
		code = elink.CodeErrProudctFeatureUndefined
	}
}

// 添加或删除通用设备
func (this *DevicesController) addDelGernalDevices(isDel bool, pid int) int {
	var sn []string
	var isArray bool
	var err error

	sns := jsoniter.Get(this.Input.Payload, "payload", "sn")
	switch sns.ValueType() {
	case jsoniter.StringValue:
		sn = append(sn, sns.ToString())
	case jsoniter.ArrayValue:
		isArray = true
		req := &DevMultiSnRequest{}
		if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
			return elink.CodeErrSysInvalidParameter
		}
		sn = req.Payload.Sn
	default:
		return elink.CodeErrSysInvalidParameter
	}
	if len(sn) == 0 {
		return elink.CodeErrSysInvalidParameter
	}

	snSuc := []string{}
	// 处理要添加或删除的设备
	for _, v := range sn {
		if models.HasGeneralDevice(pid, v) { // 设备存在
			if isDel {
				if err = models.DeleteGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		} else { // 设备不存在
			if !isDel {
				if err = models.CreateGeneralDevice(pid, v); err != nil {
					logs.Debug(err)
					continue
				}
			}
		}
		snSuc = append(snSuc, v)
	}
	if len(snSuc) == 0 {
		return elink.CodeErrDeviceCommandOperationFailed
	}

	var py interface{}
	if isArray {
		py = &DevMultiSnPy{pid, snSuc}
	} else {
		if snSuc[0] == "" {
			return elink.CodeErrDeviceCommandOperationFailed
		}
		py = &DevSnPy{pid, snSuc[0]}
	}

	err = this.WriteResponsePyServerJSON(elink.CodeSuccess, py)
	if err != nil {
		return elink.CodeErrSysException
	}
	return elink.CodeSuccess
}
