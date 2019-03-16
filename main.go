package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gogate/npis"

	"github.com/astaxie/beego"

	_ "github.com/slzm40/gogate/models/pdtModels"
	_ "github.com/slzm40/gogate/routers"
	_ "github.com/slzm40/gomo/misc"
)

func main() {
	logs.EnableFuncCallDepth(false)
	if npis.NpiAppInit() != nil {
		panic("main: npi app init failed")
	}

	beego.Run()
}
