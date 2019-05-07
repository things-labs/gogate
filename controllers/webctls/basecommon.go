package webctls

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	islogin := false
	u := this.GetSession("username")
	if u != nil {
		islogin = true
		this.Data["IsLogin"] = true
	}

	if !islogin {
		this.Redirect("/login", 302)
		this.StopRun()
	}
}
