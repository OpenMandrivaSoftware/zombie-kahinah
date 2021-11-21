package controllers

import (
	"github.com/OpenMandrivaSoftware/zombie-kahinah/models"
)

type UserController struct {
	BaseController
}

func (u *UserController) Get() {
	userStr := models.IsLoggedIn(&u.Controller)
	if userStr == "" {
		u.Abort("403")
	}

	user := models.FindUser(userStr)
    _ = user

	u.TplName = ""
}
