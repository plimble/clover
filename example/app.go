package main

import (
	"github.com/plimble/ace"
	"github.com/plimble/clover"
)

type App struct {
	auth     *clover.AuthServer
	resource *clover.ResourceServer
}

func (a *App) grantScreen(c *ace.C) {
	client, scopes, resp := a.auth.ValidateAuthorize(c.Writer, c.Request)
	if resp.IsError() {
		resp.Write(c.Writer)
		return
	}

	c.HTML("authorize.html", map[string]interface{}{
		"client": client,
		"scopes": scopes,
	})
}

func (a *App) grant(c *ace.C) {
	approve := c.MustPostString("approve", "")
	resp := a.auth.Authorize(c.Writer, c.Request, approve == "approve")
	resp.Write(c.Writer)
}

func (a *App) signinScreen(c *ace.C) {
	c.HTML("signin.html", nil)
}

func (a *App) signin(c *ace.C) {
	username := c.MustPostString("username", "")
	password := c.MustPostString("password", "")
	c.Session.Set("username", username)
	c.Session.Set("password", password)

	c.Redirect(c.MustQueryString("next", ""))
}

func (a *App) token(c *ace.C) {
	resp := a.auth.Token(c.Writer, c.Request)
	resp.Write(c.Writer)
}

func (a *App) callback(c *ace.C) {
	if c.Request.FormValue("err") != "" {
		c.String(500, "%s %s", c.Request.FormValue("err"), c.Request.FormValue("desc"))
	} else {
		c.String(200, c.Request.FormValue("code"))
	}
}

func (a *App) home(c *ace.C) {
	_, resp := a.resource.VerifyAccessToken(c.Writer, c.Request, "read_my_timeline")
	if resp.IsError() {
		resp.Write(c.Writer)
		return
	}

	c.HTML("home.html", nil)
}
