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
	defer conn.Close()

	p := elinkws.NewProvider(conn)
	p.Run()
	// for {

	// 	mt, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	logs.Debug("recv: %s", message)
	// 	err = conn.WriteMessage(mt, message)
	// 	if err != nil {
	// 		logs.Debug("write:", err)
	// 		break
	// 	}
	// }
}
