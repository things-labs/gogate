package routers

import (
	"github.com/astaxie/beego"
	"github.com/slzm40/gogate/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/npi", &controllers.NpiController{})
	beego.Router("/login", &controllers.LoginController{})
}
