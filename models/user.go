package models

import (
	"net/url"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/orm"
)

type User struct {
	Id          uint64 `orm:"auto;pk"`
	Email       string `orm:"type(text)"`
	Integration string `orm:"type(text)"` // abf service user id

	ApiKey string `orm:"type(text)"`

	Permissions []*UserPermission `orm:"rel(m2m);on_delete(set_null)"`
	Karma       []*Karma          `orm:"reverse(many);on_delete(set_null)"`

	BuildLists []*BuildList `orm:"null;reverse(many);on_delete(set_null)"`
	Advisories []*Advisory  `orm:"null;reverse(many);on_delete(set_null)"`
}

func (u *User) String() string {
	return u.Email
}

func (u *User) Save() {
	o := orm.NewOrm()
	o.Update(u)
}

type UserPermission struct {
	Id         uint64  `orm:"auto;pk"`
	Permission string  `orm:"type(text);unique"`
	Users      []*User `orm:"null;reverse(many);on_delete(set_null)"`
}

func (u *UserPermission) Save() {
	o := orm.NewOrm()
	o.Update(u)
}

// Finds the user with the given email address. If a user doesn't exist,
// it creates one and returns it.
//
// This is more commonly used with persona for login details, and with
// integrators to add buildlists.
func FindUser(email string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("Email", email).One(&user)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	} else if err != nil {
		user = User{
			Email: email,
		}
		o.Insert(&user)
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func FindUserApi(apikey string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("ApiKey", apikey).One(&user)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	} else if err != nil {
		// No such User
		return nil
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func FindUserNoCreate(email string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("Email", email).One(&user)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	} else if err != nil {
		return nil
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func PermAbortCheck(c *beego.Controller, perm string) {
	user := IsLoggedIn(c)
	if user != "" {

		model := FindUser(user)
		for _, v := range model.Permissions {
			if v.Permission == perm {
				return
			}
		}

	}

	c.Ctx.Request.Form = url.Values{}
	c.Ctx.Request.Form.Set("permission", perm)
	c.Ctx.Request.Form.Set("xsrf", c.XSRFToken())
	c.Abort("550")
}

func PermCheck(c *beego.Controller, perm string) bool {
	user := IsLoggedIn(c)
	if user != "" {

		model := FindUser(user)
		for _, v := range model.Permissions {
			if v.Permission == perm {
				return true
			}
		}

	}
	return false
}

func PermRegister(perm string) {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	num, err := qt.Filter("Permission", perm).Count()
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	}

	if num > 0 {
		return
	}

	permission := UserPermission{Permission: perm}
	_, err = o.Insert(&permission)
	if err != nil {
		panic(err)
	}
}

func PermGet(perm string) *UserPermission {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	var p UserPermission
	err := qt.Filter("Permission", perm).One(&p)
	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	} else if err != nil {
		return nil
	}

	return &p
}

func PermGetAll() []*UserPermission {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	var perms []*UserPermission
	_, err := qt.All(&perms)

	if err != nil && err.Error() != orm.ErrNoRows.Error() {
		panic(err)
	}

	for _, v := range perms {
		o.LoadRelated(v, "Users")
	}

	return perms
}
