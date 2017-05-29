package controllers

type MainController struct {
	BaseController
}

// This Get() function displays the list of
// endpoints that the application has:
//
// --> See all packages queued for testing
// --> See all packages in testing
// --> See all packages queued for updates
// --> See all packages in updates
func (this *MainController) Get() {
	this.Data["Title"] = "Main"
	this.Data["Loc"] = 0
	this.TplName = "index.tpl"
}
