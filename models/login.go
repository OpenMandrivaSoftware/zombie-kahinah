package models

import (
	beego "github.com/beego/beego/v2/adapter"
	"github.com/xiam/to"
)

const (
	PERMISSION_ADMIN     = "kahinah.admin"
	PERMISSION_QA        = "kahinah.qa"
	PERMISSION_WHITELIST = "kahinah.whitelist"

	PERMISSION_ADVISORY = "kahinah.advisory"
	PERMISSION_API      = "kahinah.api"
)

func IsLoggedIn(controller *beego.Controller) string {
	// check api key header first
	apiKey := controller.Ctx.Input.Header("X-Kahinah-Key")

	if apiKey != "" {
		user := FindUserApi(apiKey)
		if user != nil {
			for _, v := range user.Permissions {
				if v.Permission == PERMISSION_API {
					return user.Email
				}
			}
		}
	}

	// check persona
	session := controller.GetSession("github")
	if session == nil {
		return ""
	}
	return to.String(session)
}
