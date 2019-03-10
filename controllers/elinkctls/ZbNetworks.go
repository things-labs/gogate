package elinkctls

import (
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"
)

type ZbNetworkCtrlController struct {
	ctrl.CtrlController
}

// 开启zigbee网络
func (this *ZbNetworkCtrlController) Post() {
	v := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &v); err != nil {
		return
	}

	// TODO: open zigbee network

	rsp := ctrl.CtrlResponse{
		PacketID: v.PacketID,
	}

	rsp.Code = elink.CodeSuccess
	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
	}

	elink.WriteResponse(this.Input.Client, this.Input.Topic, out)
	logs.Debug("open zigbee")
}

// 关闭zigbee网络
func (this *ZbNetworkCtrlController) Delete() {
	v := ctrl.CtrlRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, &v); err != nil {
		return
	}

	// TODO: close zigbee network

	rsp := ctrl.CtrlResponse{
		PacketID: v.PacketID,
	}

	rsp.Code = elink.CodeSuccess
	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
	}

	elink.WriteResponse(this.Input.Client, this.Input.Topic, out)
	logs.Debug("close zigbee")
}
