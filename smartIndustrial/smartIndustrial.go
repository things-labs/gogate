package main

import (
	"github.com/thinkgos/gogate/apps/broad"
	"github.com/thinkgos/gogate/apps/elinkmd"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/plugin/discover"
	_ "github.com/thinkgos/gogate/smartIndustrial/routers"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego"
)

func init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	// 注册设备模型初始化函数
	models.RegisterDbTableInitFunc(models.GeneralDeviceDbTableInit)
	models.RegisterDbTableInitFunc(models.ZbDeviceDbTableInit)
}

func main() {
	elink.RegisterTopicInfo(misc.Mac(), elinkmd.ProductKey) // 注册网关产品Key
	misc.CfgInit()
	misc.LogsInit()
	err := models.DbInit()
	if err != nil {
		panic(err)
	}
	broad.BroadInit()
	go discover.Run()
	beego.Run()
}
