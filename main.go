package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/slzm40/gogate/models"
	_ "github.com/slzm40/gogate/routers"
	"github.com/slzm40/gomo/misc"
)

func main() {
	//	if npis.NpiAppInit() != nil {
	//		panic("main: npi app init failed")
	//	}
	v, err := misc.APPCfg.GetValue("COM0", "Name")
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug(v)

	beego.Run()
}
