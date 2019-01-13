package controllers

import (
	"github.com/astaxie/beego"
	"github.com/json-iterator/go"
	"github.com/slzm40/gogate/npis"
)

type NpiController struct {
	beego.Controller
}

func (this *NpiController) Post() {
	if jsoniter.Valid(this.Ctx.Input.RequestBody) {
		var any jsoniter.Any = jsoniter.Get(this.Ctx.Input.RequestBody)

		cmd := any.Get("cmd").ToUint32()
		npis.DoneCmd(uint16(cmd))
	}

	this.Data["json"] = `{"status":"success"}`
	this.ServeJSON()
}
