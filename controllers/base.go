package controllers

import (
	"html/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/robxu9/zombie-kahinah/models"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	this.Data["xsrf_token"] = this.XSRFToken()
	this.Data["xsrf_data"] = template.HTML(this.XSRFFormHTML())
	this.Data["user_login"] = models.IsLoggedIn(&this.Controller)

	this.Data["copyright"] = time.Now().Year()
}
