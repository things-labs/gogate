package routers

import (
	"github.com/astaxie/beego"
	"github.com/slzm40/gogate/controllers/elinkctls"
	"github.com/slzm40/gogate/controllers/webctls"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"
)

// web router setting
func init() {
	beego.Router("/", &webctls.HomeController{})
	beego.Router("/login", &webctls.LoginController{})
}

// elink router setting
func init() {
	elink.Router(ctrl.ChannelCtrl, "devices", &elinkctls.DevicesCtrlController{})
	elink.Router(ctrl.ChannelCtrl, "zigbee.network", &elinkctls.ZbNetworkCtrlController{})
}
