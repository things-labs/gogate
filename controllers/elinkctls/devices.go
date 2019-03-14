package elinkctls

import (
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/slzm40/gogate/models/pdtModels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"
)

type DevicesCtrlController struct {
	ctrl.CtrlController
}

//// 获取设备列表
//func (this *DevicesCtrlController) Get() {

//}
type Xpayload struct {
	ProductID uint32 `json :"productID"`
	Sn        string `json:"sn"`
}

// 添加设备
func (this *DevicesCtrlController) Post() {
	dealAddOrDelDevice(false, this)
}

// 删除设备
func (this *DevicesCtrlController) Delete() {
	dealAddOrDelDevice(true, this)
}

func dealAddOrDelDevice(isDel bool, dc *DevicesCtrlController) {
	var err error
	var sn []string
	var snSuc []string = []string{}

	code := int(1)
	req := ctrl.CtrlRequest{}
	if err = jsoniter.Unmarshal(dc.Input.Payload, &req); err != nil {
		return
	}

	logs.Info(req.Payload)
	pid, ok := req.Payload["productID"].(int)
	if !ok {
		pid = 0
	}
	sn = req.Payload["sn"].([]interface{})

	logs.Info(pid, sn)
	return
	//	logs.Debug(reflect.TypeOf(v.Payload["sn"]).Kind())
	//	return
	//	any := jsoniter.Wrap(req.Payload)
	//	if err = any.LastError(); err != nil {
	//		logs.Error(err)
	//		return

	//	pid := any.Get("productID").ToInt()
	//	snAny := any.Get("sn")
	//	snType := snAny.ValueType()

	//	if snType == jsoniter.ArrayValue {
	//		vals := reflect.ValueOf(snAny.GetInterface())
	//		logs.Info(vals.Kind())
	//		for i := 0; i < vals.Len(); i++ {
	//			logs.Info(vals.Index(i).Kind())
	//			s := reflect.ValueOf(vals.Index(i)).String()
	//			logs.Info(s)
	//			sn = append(sn, s)
	//			return
	//		}
	//	}

	if pid == 0 || len(sn) == 0 {
		code = elink.CodeErrSysInvalidParameter
	} else if !pdtModels.HasProduct(pid) {
		code = 200
	} else {
		for _, v := range sn {
			if !pdtModels.HasDevice(v, pid) {
				if isDel {
					err = pdtModels.DeleteDevice(v, pid)
				} else {
					err = pdtModels.CreateDevice(v, pid)
				}
				if err != nil {
					logs.Debug(err)
					continue
				}
				snSuc = append(snSuc, v)
			}
			code = elink.CodeSuccess
		}
	}
	rsp := ctrl.CtrlResponse{
		PacketID: req.PacketID,
	}

	rsp.Code = code
	if rsp.Code != elink.CodeSuccess {
		msg := elink.CodeErrorMessage(code)
		rsp.CodeDetail = msg.Detail
		rsp.Message = msg.Message
	} else {
		rsp.Payload = make(map[string]interface{})
		rsp.Payload["productID"] = pid
		rsp.Payload["sn"] = snSuc
	}

	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
	}

	elink.WriteResponse(dc.Input.Client, dc.Input.Topic, out)
}
