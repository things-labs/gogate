package main

import (
	"github.com/thinkgos/gogate/apps/mq"
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/plugin/discover"

	"github.com/astaxie/beego"

	_ "github.com/thinkgos/gogate/models"
	_ "github.com/thinkgos/gogate/smartIndustrial/routers"
)

func main() {
	misc.LogsInit()
	mq.MqInit()
	if npis.OpenZbApp() != nil {
		panic("main: npi app init failed")
	}
	go discover.Run()
	beego.Run()
}
