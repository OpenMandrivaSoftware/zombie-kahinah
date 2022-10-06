package models

import (
	"net/url"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	uuid "github.com/satori/go.uuid"
	"github.com/xiam/to"
)

var (
	outwardUrl       = beego.AppConfig.String("outwardloc")
	parsedOutwardUrl = mustParseUrl(outwardUrl)
	githubKey        = beego.AppConfig.String("auth::githubKey")
	githubSecret     = beego.AppConfig.String("auth::githubSecret")
)

func init() {
	goth.UseProviders(github.New(githubKey, githubSecret, outwardUrl+"/auth/login/callback", "user:email"))
}

type GithubLogoutController struct {
	beego.Controller
}

func (this *GithubLogoutController) Get() {
	this.DestroySession()
	referrer := this.Ctx.Request.Referer()
	if referrer == "" {
		this.Redirect(outwardUrl, 307)
	} else {
		this.Redirect(referrer, 307)
	}
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
		urlParsed, err := url.Parse(referrer)
		if err != nil {
			this.SetSession("login-referrer", referrer)
		} else {
			if urlParsed.Host != parsedOutwardUrl.Host {
				// Maybe a hostname that's wrong? Use our own.
				urlParsed.Host = parsedOutwardUrl.Host
				this.SetSession("login-referrer", urlParsed.String())
			} else {
				this.SetSession("login-referrer", referrer)
			}
		}
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

func mustParseUrl(urlStr string) *url.URL {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return parsed
}
