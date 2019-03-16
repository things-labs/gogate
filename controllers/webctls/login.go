package webctls

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	logs.Debug(this.Ctx.Input.Param(":id"))
	//this.TplName = "login.html"
}

func (this *LoginController) Post() {
	// logs.Debug(this.Ctx.Request.Body)
	logs.Debug(this.Data["username"])
	logs.Debug(this.Data["password"])
	this.Redirect("/", 302)
}

// 生成密钥 加盐法
func generatePwd(username, password string) string {
	h := md5.New()
	io.WriteString(h, password)
	md5Pwd := fmt.Sprintf("%x", h.Sum(nil))

	// 加盐值加密
	// 用户名+`@#$%`+md5Pwd+`^&*()`拼接
	io.WriteString(h, username)
	io.WriteString(h, `@#$%`)
	io.WriteString(h, md5Pwd)
	io.WriteString(h, `^&*()`)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkPassword(pwd string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", pwd); !ok {
		return false
	}
	return true
}

func checkUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", username); !ok {
		return false
	}
	return true
}
