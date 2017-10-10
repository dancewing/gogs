// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/api"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
	"github.com/gogits/gogs/pkg/setting"
)

func Search(c *context.APIContext) {

}

func GetInfo(c *context.APIContext) {
	u, err := models.GetUserByName(c.Params(":username"))
	if err != nil {
		if errors.IsUserNotExist(err) {
			c.Status(404)
		} else {
			c.Error(500, "GetUserByName", err)
		}
		return
	}

	// Hide user e-mail when API caller isn't signed in.
	if !c.IsLogged {
		u.Email = ""
	}
	c.JSON(200, api.ConvertUserToAPI(u))
}

func afterLogin(c *context.Context, u *models.User, remember bool) {
	if remember {
		days := 86400 * setting.LoginRememberDays
		c.SetCookie(setting.CookieUserName, u.Name, days, setting.AppSubURL, "", setting.CookieSecure, true)
		c.SetSuperSecureCookie(u.Rands+u.Passwd, setting.CookieRememberName, u.Name, days, setting.AppSubURL, "", setting.CookieSecure, true)
	}

	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.Name)
	c.Session.Delete("twoFactorRemember")
	c.Session.Delete("twoFactorUserID")

	// Clear whatever CSRF has right now, force to generate a new one
	c.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubURL)
	if setting.EnableLoginStatusCookie {
		c.SetCookie(setting.LoginStatusCookieName, "true", 0, setting.AppSubURL)
	}
	//redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to"))
	//c.SetCookie("redirect_to", "", -1, setting.AppSubURL)

}

func AuthenticatePost(c *context.Context, f form.SignIn) {

	u, err := models.UserSignIn(f.UserName, f.Password)
	if err != nil {
		if errors.IsUserNotExist(err) {
			c.JSON(500, nil)
		} else {
			c.ServerError("UserSignIn", err)
		}
		return
	}

	if !u.IsEnabledTwoFactor() {
		afterLogin(c, u, f.Remember)
		return
	}

	c.Session.Set("twoFactorRemember", f.Remember)
	c.Session.Set("twoFactorUserID", u.ID)
	c.Redirect(setting.AppSubURL + "/user/login/two_factor")
}

func AuthenticateGet(c *context.APIContext) {
	if c.IsLogged {
		c.JSONString(c.User.FullName)
	} else {
		c.JSONString("")
	}
}

func GetAccount(c *context.APIContext) {
	if c.IsLogged {
		user := api.ConvertUserToAPI(c.User)
		c.JSONSuccess(user)
	} else {
		c.JSON(500, nil)
	}
}
