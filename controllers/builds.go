package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"gitea.tsn.sh/robert/zombie-kahinah/util"

	"gitea.tsn.sh/robert/zombie-kahinah/integration"
	"gitea.tsn.sh/robert/zombie-kahinah/models"
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/xiam/to"
)

const (
	block_karma = 9999
	push_karma  = 9999
)

var (
	maintainer_karma = to.Int64(beego.AppConfig.String("karma::maintainerkarma"))
	maintainer_hours = to.Int64(beego.AppConfig.String("karma::maintainerhours"))
)

type ByUpdateDate []*models.BuildList

func (b ByUpdateDate) Len() int {
	return len(b)
}
func (b ByUpdateDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByUpdateDate) Less(i, j int) bool {
	return b[i].Updated.Unix() > b[j].Updated.Unix()
}

//
// --------------------------------------------------------------------
// LISTS
// --------------------------------------------------------------------
//

// ALL BUILDS
type BuildsController struct {
	BaseController
}

func (this *BuildsController) Get() {
	pageint, err := this.GetInt("page")
	if err != nil {
		pageint = 1
	} else if pageint <= 0 {
		pageint = 1
	}

	page := int64(pageint)

	var packages []*models.BuildList

	o := orm.NewOrm()

	qt := o.QueryTable(new(models.BuildList))

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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").All(&packages)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Packages")
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByUpdateDate(packages))

	this.Data["Title"] = "Builds"
	this.Data["Loc"] = 1
	this.Data["Tab"] = 4
	this.Data["Packages"] = packages
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplName = "builds/builds_list.tpl"
}

// REJECTED BUILDS
type RejectedController struct {
	BaseController
}

func (this *RejectedController) Get() {
	pageint, err := this.GetInt("page")
	if err != nil {
		pageint = 1
	} else if pageint <= 0 {
		pageint = 1
	}

	page := int64(pageint)

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	cnt, err := qt.Filter("status", models.STATUS_REJECTED).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.STATUS_REJECTED).All(&packages)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Packages")
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByUpdateDate(packages))

	this.Data["Title"] = "Rejected"
	this.Data["Loc"] = 1
	this.Data["Tab"] = 3
	this.Data["Packages"] = packages
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplName = "builds/builds_list.tpl"
}

// PUBLISHED BUILDS
type PublishedController struct {
	BaseController
}

func (this *PublishedController) Get() {
	filterPlatform := this.GetString("platform")

	pageint, err := this.GetInt("page")
	if err != nil {
		pageint = 1
	} else if pageint <= 0 {
		pageint = 1
	}

	page := int64(pageint)

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	if filterPlatform != "" {
		qt = qt.Filter("platform", filterPlatform)
	}

	cnt, err := qt.Filter("status", models.STATUS_PUBLISHED).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.STATUS_PUBLISHED).All(&packages)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Packages")
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByUpdateDate(packages))

	this.Data["Title"] = "Published"
	this.Data["Loc"] = 1
	this.Data["Tab"] = 2
	this.Data["Packages"] = packages
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplName = "builds/builds_list.tpl"
}

// TESTING BUILDS
type ByBuildDate []*models.BuildList

func (b ByBuildDate) Len() int {
	return len(b)
}
func (b ByBuildDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByBuildDate) Less(i, j int) bool {
	return b[i].BuildDate.Unix() > b[j].BuildDate.Unix()
}

type TestingController struct {
	BaseController
}

func (this *TestingController) Get() {
	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.STATUS_TESTING).All(&packages)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		log.Println(err)
		this.Abort("500")
	}

	pkgkarma := make(map[string]string)

	for _, v := range packages {
		totalKarma := getTotalKarma(v.Id)

		pkgkarma[to.String(v.Id)] = to.String(totalKarma)

		o.LoadRelated(v, "Packages")
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByBuildDate(packages))

	if models.IsLoggedIn(&this.Controller) != "" {
		this.Data["LoggedIn"] = true
	}

	if models.PermCheck(&this.Controller, models.PERMISSION_QA) {
		this.Data["QAControls"] = true
	}

	this.Data["Title"] = "Testing"
	this.Data["Loc"] = 1
	this.Data["Tab"] = 1
	this.Data["Packages"] = packages
	this.Data["PkgKarma"] = pkgkarma
	this.Data["Entries"] = num
	this.TplName = "builds/generic_list.tpl"
}

//
// --------------------------------------------------------------------
// INDIVIDUAL BUILD LOOK
// --------------------------------------------------------------------
//

type BuildController struct {
	BaseController
}

func (this *BuildController) Get() {
	id := to.Uint64(this.Ctx.Input.Param(":id"))

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err != nil && err.Error() == orm.ErrNoRows.Error() {
		this.Abort("404")
	} else if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	if pkg.Changelog != "" {
		resp, err := http.Get(pkg.Changelog)
		if err == nil {
			defer resp.Body.Close()
			changelog, _ := ioutil.ReadAll(resp.Body)
			this.Data["Changelog"] = this.processChangelog(string(changelog))
		} else {
			this.Data["Changelog"] = "Failed to retrieve changelog: " + err.Error()
		}
	}

	this.Data["Commits"] = integration.Commits(&pkg)

	this.Data["Url"] = integration.Url(&pkg)

	o.LoadRelated(&pkg, "Submitter")
	o.LoadRelated(&pkg, "Packages")

	// karma controls
	totalKarma := getTotalKarma(id) // get total karma

	votes := make([]util.Pair, 0) // *models.Karma, int

	// load karma totals
	var inOrder []*models.Karma
	kt := o.QueryTable(new(models.Karma))
	kt.Filter("List__Id", id).OrderBy("Time").All(&inOrder)

	// only count most recent votes
	for _, v := range inOrder {
		o.LoadRelated(v, "User")

		pair := util.Pair{}
		pair.Key = v

		switch v.Vote {
		case models.KARMA_UP:
			pair.Value = 1
		case models.KARMA_DOWN:
			pair.Value = 2
		case models.KARMA_MAINTAINER:
			pair.Value = 1
		case models.KARMA_BLOCK:
			pair.Value = 2
		case models.KARMA_PUSH:
			pair.Value = 1
		case models.KARMA_NONE:
			if v.Comment != "" {
				pair.Value = 0
			} else {
				continue // no karma and no comment? useless
			}
		case models.KARMA_FINALIZE:
			pair.Value = 3
		case models.KARMA_CLEAR:
			pair.Value = 4
		}

		votes = append(votes, pair)
	}

	this.Data["Votes"] = votes
	this.Data["Karma"] = totalKarma

	this.Data["UserVote"] = 0
	user := models.IsLoggedIn(&this.Controller)
	if user != "" {
		kt := o.QueryTable(new(models.Karma))
		var userkarma models.Karma
		err = kt.Filter("User__Email", user).Filter("List__Id", id).OrderBy("-Time").Limit(1).One(&userkarma)
		if err != nil && err.Error() != orm.ErrNoRows.Error() {
			log.Println(err)
		} else if err == nil {
			switch userkarma.Vote {
			case models.KARMA_UP:
				this.Data["UserVote"] = 1
			case models.KARMA_MAINTAINER:
				this.Data["UserVote"] = 2
			case models.KARMA_DOWN:
				this.Data["UserVote"] = -1
			}

			this.Data["KarmaCommentPrev"] = userkarma.Comment
		}

		if models.PermCheck(&this.Controller, models.PERMISSION_QA) {
			this.Data["QAControls"] = true
		}

		upthreshold, err := beego.AppConfig.Int64("karma::upperkarma")
		if err != nil {
			panic(err)
		}

		downthreshold, err := beego.AppConfig.Int64("karma::lowerkarma")
		if err != nil {
			panic(err)
		}

		if totalKarma >= int(upthreshold) || totalKarma <= int(downthreshold) {
			this.Data["FinalizeControls"] = true
		}

	}

	// karma controls end

	this.Data["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	this.Data["Loc"] = 1
	if pkg.Status == models.STATUS_TESTING {
		this.Data["Tab"] = 1
		if user != "" {
			this.Data["KarmaControls"] = true
			if pkg.Submitter != nil && pkg.Submitter.Email == user {
				this.Data["MaintainerControls"] = true
				this.Data["MaintainerHoursNeeded"] = maintainer_hours
				if time.Since(pkg.BuildDate).Hours() >= float64(maintainer_hours) {
					this.Data["MaintainerTime"] = true
					delete(this.Data, "MaintainerHoursNeeded")
				}
			}
		}
	} else if pkg.Status == models.STATUS_PUBLISHED {
		this.Data["Tab"] = 2
	} else if pkg.Status == models.STATUS_REJECTED {
		this.Data["Tab"] = 3
	} else {
		this.Data["Tab"] = 4
	}
	this.Data["Package"] = pkg
	this.Data["SourceEVR"] = pkg.SourceEVR()
	this.TplName = "builds/build.tpl"
}

func (this *BuildController) Post() {
	id := to.Uint64(this.Ctx.Input.Param(":id"))

	postType := this.GetString("type")
	if postType != "Neutral" && postType != "Up" && postType != "Down" && postType != "Maintainer" && postType != "QABlock" && postType != "QAPush" && postType != "QAClear" && postType != "Finalize" {
		this.Abort("400")
	}

	comment := this.GetString("comment")

	user := models.IsLoggedIn(&this.Controller)
	if user == "" {
		this.Abort("403") // MUST be logged in
	}

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err != nil && err.Error() == orm.ErrNoRows.Error() {
		this.Abort("404")
	} else if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	o.LoadRelated(&pkg, "Submitter")

	if postType == "Maintainer" {
		if pkg.Submitter.Email != user {
			this.Abort("403")
		} else {
			if time.Since(pkg.BuildDate).Hours() < float64(maintainer_hours) { // week
				this.Abort("400")
			}
		}
	} else if postType == "QABlock" || postType == "QAPush" || postType == "QAClear" {
		models.PermAbortCheck(&this.Controller, models.PERMISSION_QA)
	} else {
		// whitelist stuff
		if Whitelist {
			perm := models.PermCheck(&this.Controller, models.PERMISSION_WHITELIST)
			if !perm {
				flash := beego.NewFlash()
				flash.Warning("Sorry, the whitelist is on and you are not allowed to vote.")
				flash.Store(&this.Controller)
				this.Get()
				return
			}
		}
	}

	if postType == "Finalize" {
		karmaTotal := getTotalKarma(id)

		upthreshold, err := beego.AppConfig.Int64("karma::upperkarma")
		if err != nil {
			panic(err)
		}

		downthreshold, err := beego.AppConfig.Int64("karma::lowerkarma")
		if err != nil {
			panic(err)
		}

		if karmaTotal >= int(upthreshold) {
			pkg.Status = models.STATUS_PUBLISHED
			o.Update(&pkg)
			go integration.Publish(&pkg)
		} else if karmaTotal <= int(downthreshold) {
			pkg.Status = models.STATUS_REJECTED
			o.Update(&pkg)
			go integration.Reject(&pkg)
		} else {
			this.Abort("400") // bad request - we can't finalize!
		}

		go util.MailModel(&pkg)
	}

	var userkarma models.Karma

	userkarma.List = &pkg
	userkarma.User = models.FindUser(user)
	if postType == "Up" {
		userkarma.Vote = models.KARMA_UP
	} else if postType == "Maintainer" {
		userkarma.Vote = models.KARMA_MAINTAINER
	} else if postType == "QABlock" {
		userkarma.Vote = models.KARMA_BLOCK
	} else if postType == "QAPush" {
		userkarma.Vote = models.KARMA_PUSH
	} else if postType == "Neutral" {
		userkarma.Vote = models.KARMA_NONE
	} else if postType == "Down" {
		userkarma.Vote = models.KARMA_DOWN
	} else if postType == "QAClear" {
		userkarma.Vote = models.KARMA_CLEAR
	} else {
		userkarma.Vote = models.KARMA_FINALIZE
	}

	userkarma.Comment = comment
	o.Insert(&userkarma)

	if postType == "QAClear" {
		pkg.Status = models.STATUS_REJECTED
		o.Update(&pkg)
	}

	this.Get()
}

func getTotalKarma(id uint64) int {
	o := orm.NewOrm()
	kt := o.QueryTable(new(models.Karma))

	var karma []*models.Karma
	kt.Filter("List__Id", id).OrderBy("-Time").All(&karma)

	set := util.NewSet()
	totalKarma := 0

	// only count most recent votes
	for _, v := range karma {
		o.LoadRelated(v, "User")

		if set.Contains(v.User.Email) {
			continue // we've already counted this person's most recent vote
		}

		switch v.Vote {
		case models.KARMA_UP:
			totalKarma++
		case models.KARMA_DOWN:
			totalKarma--
		case models.KARMA_MAINTAINER:
			totalKarma += int(maintainer_karma)
		case models.KARMA_BLOCK:
			totalKarma -= int(block_karma)
		case models.KARMA_PUSH:
			totalKarma += int(push_karma)
		case models.KARMA_CLEAR:
			fallthrough
		case models.KARMA_FINALIZE:
			continue // we ignore these kinds of karma
		}

		set.Add(v.User.Email)
	}

	return totalKarma
}

func (this *BuildController) processChangelog(changelog string) string {
	toreturn := ""
	open := true
	for _, c := range changelog {
		if c == '<' {
			toreturn += string(c)
			toreturn += "email hidden"
			open = false
		} else if c == '>' {
			toreturn += string(c)
			open = true
		} else if open {
			toreturn += string(c)
		}
	}
	return toreturn
}
