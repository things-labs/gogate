package main

import (
	"github.com/slzm40/gogate/apps/npis"

	"github.com/astaxie/beego"

	_ "github.com/slzm40/gogate/apps/mq"
	_ "github.com/slzm40/gogate/models/devmodels"
	_ "github.com/slzm40/gogate/routers"
	_ "github.com/slzm40/gomo/misc"
)

func main() {
	if npis.ZbAppInit() != nil {
		panic("main: npi app init failed")
	}

	beego.Run()
}
