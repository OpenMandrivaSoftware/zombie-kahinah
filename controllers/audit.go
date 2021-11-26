package controllers

import (
	"log"

	"gitea.tsn.sh/robert/zombie-kahinah/models"
	"github.com/beego/beego/v2/adapter/orm"
)

type AuditController struct {
	BaseController
}

func (this *AuditController) Get() {
	pageint, err := this.GetInt("page")
	if err != nil {
		pageint = 1
	} else if pageint <= 0 {
		pageint = 1
	}

	page := int64(pageint)

	var karma []*models.Karma

	o := orm.NewOrm()

	qt := o.QueryTable(new(models.Karma))

	cnt, err := qt.Count()
	if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	totalpages := cnt / 50
	if cnt%50 != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Time").All(&karma)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range karma {
		o.LoadRelated(v, "List")
		o.LoadRelated(v.List, "Packages")
		o.LoadRelated(v, "User")
	}

	this.Data["Title"] = "Audit Log"
	this.Data["Loc"] = 2
	this.Data["Karma"] = karma
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplName = "audit/audit_list.tpl"
}
