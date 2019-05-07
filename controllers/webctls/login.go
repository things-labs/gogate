package webctls

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	uname := this.Input().Get("username")
	pwd := this.Input().Get("password")
	if len(uname) == 0 || len(pwd) == 0 {
		this.Redirect("/login.html", 302)
		return
	}

	newPwd := GeneratePwd(uname, pwd)
	lname := beego.AppConfig.String("username")
	lpwd := beego.AppConfig.String("password")
	if strings.EqualFold(uname, lname) && strings.EqualFold(newPwd, lpwd) {
		this.SetSession("username", uname)
		this.Redirect("/index.html", 302) // 重定向到首页
	}

	this.Redirect("/login.html", 302)
}

// GeneratePwd 生成密钥 加盐法 用户名+`@#$%`+md5Pwd+`^&*()`拼接
func GeneratePwd(username, password string) string {
	h := md5.New()
	io.WriteString(h, password)
	md5Pwd := fmt.Sprintf("%x", h.Sum(nil))
	// 加盐值加密
	io.WriteString(h, username)
	io.WriteString(h, `@#$%`)
	io.WriteString(h, md5Pwd)
	io.WriteString(h, `^&*()`)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// func checkPassword(pwd string) bool {
// 	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", pwd); !ok {
// 		return false
// 	}
// 	return true
// }

// func checkUsername(username string) (b bool) {
// 	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", username); !ok {
// 		return false
// 	}
// 	return true
// }
