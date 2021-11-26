package models

import (
	beego "github.com/beego/beego/v2/adapter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	uuid "github.com/satori/go.uuid"
	"github.com/xiam/to"
)

var (
	outwardUrl   = beego.AppConfig.String("outwardloc")
	githubKey    = beego.AppConfig.String("auth::githubKey")
	githubSecret = beego.AppConfig.String("auth::githubSecret")
)

func init() {
	goth.UseProviders(github.New(githubKey, githubSecret, outwardUrl+"/auth/login/callback", "user"))
}

type GithubCheckController struct {
	beego.Controller
}

func (this *GithubCheckController) Get() {
	session := this.GetSession("github")
	if session != nil {
		emailLogin := to.String(session)
		this.Ctx.WriteString(emailLogin)
	} else {
		this.Ctx.WriteString("")
	}
}

type GithubLogoutController struct {
	beego.Controller
}

func (this *GithubLogoutController) Get() {
	this.DestroySession()
	this.Ctx.WriteString("OK")
}

type GithubLoginController struct {
	beego.Controller
}

func (this *GithubLoginController) Get() {
	provider, err := goth.GetProvider("github")
	if err != nil {
		panic(err)
	}

	stateUUID := uuid.NewV4()
	this.SetSession("oauth2State", stateUUID.String())

	oauth2Sess, err := provider.BeginAuth(stateUUID.String())
	if err != nil {
		panic(err)
	}

	oauth2URL, err := oauth2Sess.GetAuthURL()
	if err != nil {
		panic(err)
	}

	this.SetSession("oauth2Sess", oauth2Sess.Marshal())

	referrer := this.Ctx.Request.Referer()
	if referrer == "" {
		this.SetSession("login-referrer", outwardUrl)
	} else {
		this.SetSession("login-referrer", referrer)
	}

	this.Redirect(oauth2URL, 307)
}

type GithubLoginCallbackController struct {
	beego.Controller
}

func (this *GithubLoginCallbackController) Get() {

	provider, err := goth.GetProvider("github")
	if err != nil {
		panic(err)
	}

	oauth2Sess, err := provider.UnmarshalSession(to.String(this.GetSession("oauth2Sess")))
	if err != nil {
		panic(err)
	}

	_, err = oauth2Sess.Authorize(provider, this.Ctx.Request.URL.Query())
	if err != nil {
		panic(err)
	}

	oauth2User, err := provider.FetchUser(oauth2Sess)
	if err != nil {
		panic(err)
	}

	go FindUser(oauth2User.Email)

	this.SetSession("github", oauth2User.Email)

	referrer := to.String(this.GetSession("login-referrer"))
	this.Redirect(referrer, 307)
}
