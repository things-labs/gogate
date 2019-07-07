package webctls

import (
	"github.com/thinkgos/gogate/apps/broad"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ElwsController struct {
	beego.Controller
}

func (this *ElwsController) ConnectWs() {
	if err := broad.WsHub.UpgradeWithRun(this.Ctx.ResponseWriter, this.Ctx.Request); err != nil {
		logs.Error(err)
	}
	return
}
