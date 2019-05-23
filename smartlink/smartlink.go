package main

import (
	_ "github.com/thinkgos/gogate/smartlink/routers"

	"github.com/astaxie/beego"
	_ "github.com/thinkgos/gogate/apps/mbapp"
)

func init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
}

func main() {
	// misc.CfgInit()
	// misc.LogsInit()
	// err := models.DbInit()
	// if err != nil {
	// 	panic(err)
	// }

	beego.Run()
}
