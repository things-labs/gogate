package elinkctls

import (
	"github.com/thinkgos/elink"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/memlog"

	jsoniter "github.com/json-iterator/go"
)

// DevSn 单设备负载
type DevSn struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

// DevMultiSn 多设备负载
type DevMultiSn struct {
	ProductID int      `json:"productID"`
	Sn        []string `json:"sn"`
}

// DevMultiSnRequest 多设备请求
type DevMultiSnRequest struct {
	ctrl.BaseRequest
	Payload DevMultiSn `json:"payload,omitempty"`
}

// DevicesController 设备控制器
type DevicesController struct {
	ctrl.Controller
}

// Get 获取产品Id下的设备列表
func (this *DevicesController) Get() {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	pid, err := this.AcquireParamPid()
	if err != nil {
		code = elink.CodeErrSysResourceNotSupport
		return
	}

	pInfo, err := models.LookupProduct(pid)
	if err != nil {
		code = ctrl.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case models.PTypesGeneral: // 获取通用设备
		err = this.getGernalDevices(pid)
		if err != nil {
			code = elink.CodeErrSysException
		}
	default:
		code = ctrl.CodeErrProudctFeatureUndefined
	}
}

// Post 添加设备
func (this *DevicesController) Post() {
	this.dealAddDelGernalDevices(false)
}

// Delete 删除设备
func (this *DevicesController) Delete() {
	this.dealAddDelGernalDevices(true)
}

// 获取通用设备列表
func (this *DevicesController) getGernalDevices(pid int) error {
	devs := models.FindGeneralDevice(pid)
	sns := make([]string, 0, len(devs))
	for _, v := range devs {
		sns = append(sns, v.Sn)
	}
	return this.WriteResponsePyServerJSON(elink.CodeSuccess, &DevMultiSn{pid, sns})
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
		code = ctrl.CodeErrProudctUndefined
		return
	}

	// 根据不同的设备类型分发
	switch pInfo.Types {
	case models.PTypesGeneral: // 通用设备处理s
		code = this.addDelGernalDevices(isDel, int(pid))
	default:
		code = ctrl.CodeErrProudctFeatureUndefined
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

	sucSn := []string{}
	// 处理要添加或删除的设备
	for _, v := range sn {
		if models.HasGeneralDevice(pid, v) { // 设备存在
			if isDel {
				if err = models.DeleteGeneralDevice(pid, v); err != nil {
					memlog.Debug(err)
					continue
				}
			}
		} else { // 设备不存在
			if !isDel {
				if err = models.CreateGeneralDevice(pid, v); err != nil {
					memlog.Debug(err)
					continue
				}
			}
		}
		sucSn = append(sucSn, v)
	}
	if len(sucSn) == 0 {
		return ctrl.CodeErrDeviceCommandOperationFailed
	}

	var py interface{}
	if isArray {
		py = &DevMultiSn{pid, sucSn}
	} else {
		if sucSn[0] == "" {
			return ctrl.CodeErrDeviceCommandOperationFailed
		}
		py = &DevSn{pid, sucSn[0]}
	}

	err = this.WriteResponsePyServerJSON(elink.CodeSuccess, py)
	if err != nil {
		return elink.CodeErrSysException
	}
	return elink.CodeSuccess
}
