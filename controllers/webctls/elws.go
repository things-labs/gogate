package webctls

import (
	"github.com/thinkgos/gogate/protocol/elinkws"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

type ElwsController struct {
	beego.Controller
}

func (this *ElwsController) Get() {
	conn, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		logs.Error(err)
		return
	}

	elinkws.NewProvider(conn).Run()
}
