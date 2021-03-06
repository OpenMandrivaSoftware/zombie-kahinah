package controllers

import (
	"fmt"
	"log"
	"strings"

	"gitea.tsn.sh/robert/zombie-kahinah/models"
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/orm"
)

var (
	enabledPlatforms = make(map[string]string) // [platform]PREFIX
	types            = make([]string, 4)
)

func init() {
	configPlatforms := strings.Split(beego.AppConfig.String("advisory::platforms"), ";")
	for _, v := range configPlatforms {
		parts := strings.Split(v, ":")
		enabledPlatforms[parts[0]] = parts[1]
	}

	types[0] = "Security"
	types[1] = "Recommended"
	types[2] = "Bug Fix"
	types[3] = "New Release"
}

//
// our base controller
//

type AdvisoryBaseController struct {
	BaseController
}

func (this *AdvisoryBaseController) Prepare() {
	this.BaseController.Prepare()
	this.Data["Loc"] = 2
}

//
// main controller
// shows recent advisories for enabled platforms
//

type AdvisoryMainController struct {
	AdvisoryBaseController
}

func (this *AdvisoryMainController) Get() {
	platforms := make(map[string][]*models.Advisory)

	o := orm.NewOrm()

	for k, _ := range enabledPlatforms {
		qt := o.QueryTable(new(models.Advisory)).Filter("Platform", k).Exclude("AdvisoryId", 0).OrderBy("-Issued").Limit(5)

		var advisories []*models.Advisory

		_, err := qt.All(&advisories)
		if err != nil && err.Error() != orm.ErrNoRows.Error() {
			log.Printf("error occured trying to get advisories: %s", err)
			this.Abort("500")
		}
		platforms[k] = advisories
	}

	this.Data["Tab"] = 1
	this.Data["Title"] = "Advisories"
	this.TplName = "advisories/main.tpl"

	this.Data["Platforms"] = platforms
}

//
// new controller
// shows new input
//

type AdvisoryNewController struct {
	AdvisoryBaseController
}

func (this *AdvisoryNewController) Get() {
	models.PermAbortCheck(&this.Controller, models.PERMISSION_ADVISORY)

	this.Data["FailPlatform"] = ""

	this.Data["Platforms"] = enabledPlatforms
	this.Data["Types"] = types

	this.Data["Tab"] = -1
	this.Data["Title"] = "New Advisory"
	this.TplName = "advisories/new.tpl"
}

func (this *AdvisoryNewController) Post() {
	models.PermAbortCheck(&this.Controller, models.PERMISSION_ADVISORY)

	fmt.Printf("%s\n", this.Input())

	this.Abort("500")
}
