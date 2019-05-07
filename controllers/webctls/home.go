package webctls

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.TplName = "index.html"
}
