package main

import (
	"log"
	"net/http"

	_ "github.com/things-labs/fwu"
	"github.com/thinkgos/gogate/apps/broad"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/apps/elinkmd"
	"github.com/thinkgos/gogate/apps/npis"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/plugin/discover"
	_ "github.com/thinkgos/gogate/routers"
	"github.com/thinkgos/memlog"
)

func init() {
	// 注册用户模型初始化函数
	models.RegisterDbTableInitFunc(models.UserDbTableInit)
	// 注册设备模型初始化函数
	models.RegisterDbTableInitFunc(models.GeneralDeviceDbTableInit)
	models.RegisterDbTableInitFunc(models.ZbDeviceDbTableInit)
}

func main() {
	memlog.SetLogger(memlog.AdapterConsole)
	ctrl.RegisterTopicInfo(misc.Mac(), elinkmd.ProductKey) // 注册网关产品Key
	misc.ConfigInit()
	err := models.DbInit()
	if err != nil {
		panic(err)
	}
	broad.BroadInit()

	if err = npis.OpenZbApp(); err != nil {
		panic(err)
	}
	go discover.Run()
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Printf("http listen and serve failed, %v", err)
	}
}
