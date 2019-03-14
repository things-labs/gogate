package main

import (
	"github.com/astaxie/beego"
	_ "github.com/slzm40/gogate/models/pdtModels"
	_ "github.com/slzm40/gogate/routers"
	_ "github.com/slzm40/gomo/misc"
)

func main() {
	//if npis.NpiAppInit() != nil {
	//		panic("main: npi app init failed")
	//	}

	beego.Run()
}
