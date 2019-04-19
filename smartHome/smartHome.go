package main

import (
	"github.com/thinkgos/gogate/apps/mq"
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/plugin/discover"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"

	_ "github.com/thinkgos/gogate/models"
	_ "github.com/thinkgos/gogate/smartHome/routers"
)

func main() {
	elink.RegisterTopicInfo(misc.Mac(), elinkmd.ProductKey) // 注册网关产品Key
	misc.LogsInit()
	mq.MqInit(elinkmd.ProductKey, misc.Mac())
	if npis.OpenZbApp() != nil {
		panic("main: npi app init failed")
	}

	discover.Run()
}
