package webctls

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["Website"] = "www.lchtime.com"
	this.Data["Email"] = "jgb40@qq.com"
	this.TplName = "index.html"
}
