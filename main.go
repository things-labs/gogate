package main

import (
	"github.com/thinkgos/gogate/apps/npis"

	"github.com/astaxie/beego"

	_ "github.com/thinkgos/gogate/apps/mq"
	_ "github.com/thinkgos/gogate/models/devmodels"
	_ "github.com/thinkgos/gogate/routers"
	_ "github.com/thinkgos/gomo/misc"
)

func main() {
	if npis.ZbAppInit() != nil {
		panic("main: npi app init failed")
	}

	beego.Run()
}
