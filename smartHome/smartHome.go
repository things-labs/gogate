package main

import (
	"github.com/thinkgos/gogate/apps/mq"
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/plugin/discover"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"

	"github.com/thinkgos/gogate/models"
	_ "github.com/thinkgos/gogate/smartHome/routers"
)

func init() {
	// 注册设备模型初始化函数
	models.RegisterDbTableInitFunction(models.GeneralDeviceDbTableInit)
	models.RegisterDbTableInitFunction(models.ZbDeviceDbTableInit)
}

func main() {
	elink.RegisterTopicInfo(misc.Mac(), elinkmd.ProductKey) // 注册网关产品Key
	misc.CfgInit()
	misc.LogsInit()
	err := models.DbInit()
	if err != nil {
		panic(err)
	}

	mq.MqInit(elinkmd.ProductKey, misc.Mac())
	err = npis.OpenZbApp()
	if err != nil {
		panic(err)
	}

	discover.Run()
}
