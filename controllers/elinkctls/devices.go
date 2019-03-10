package elinkctls

import (
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
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
	v := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &v); err != nil {
		return
	}

	any := jsoniter.Wrap(v.Payload)
	if err := any.LastError(); err != nil {
		logs.Error(err)
	}

	logs.Debug(any.Get("productID").ToUint32())
	logs.Debug(any.Get("sn").ToString())
	// TODO: 添加设备

	rsp := ctrl.CtrlResponse{
		PacketID: v.PacketID,
	}

	rsp.Code = elink.CodeSuccess
	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
	}

	elink.WriteResponse(this.Input.Client, this.Input.Topic, out)

}

// 删除设备
func (this *DevicesCtrlController) Delete() {
	v := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &v); err != nil {
		return
	}
	// TODO: 删除设备
	rsp := ctrl.CtrlResponse{
		PacketID: v.PacketID,
	}

	rsp.Code = elink.CodeSuccess
	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
	}

	elink.WriteResponse(this.Input.Client, this.Input.Topic, out)
}
