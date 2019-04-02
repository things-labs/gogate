package routers

import (
	"github.com/thinkgos/gogate/controllers/elinkctls"
	"github.com/thinkgos/gogate/controllers/webctls"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/protocol/elinkch/ctrl"

	"github.com/astaxie/beego"
)

// web router setting
func init() {
	beego.Router("/", &webctls.HomeController{})
	beego.Router("/login/:id([0-9]+)", &webctls.LoginController{})
}

// elink router setting
func init() {
	elink.Router(ctrl.ChannelCtrl, "devices.@", &elinkctls.DevicesController{})
	elink.Router(ctrl.ChannelCtrl, "device.commands.@", &elinkctls.DevCommandController{})
	elink.Router(ctrl.ChannelCtrl, "device.propertys.@", &elinkctls.DevPropertysController{})
	elink.Router(ctrl.ChannelCtrl, "zigbee.network", &elinkctls.ZbNetworkController{})
}
