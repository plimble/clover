package main

import (
	"github.com/plimble/ace"
	"github.com/plimble/clover"
)

type App struct {
	authServer     *clover.AuthServer
	resourceServer *clover.ResourceServer
}

func (a *App) grantScreen(c *ace.C) {
	ad, resp := a.authServer.ValidateAuthorize(c.Writer, c.Request)
	if resp != nil {
		resp.Write(c.Writer)
		return
	}

	c.HTML("authorize.html", map[string]interface{}{
		"client": ad.Client,
		"scopes": ad.Scope,
	})
}

func (a *App) grant(c *ace.C) {
	session := c.Sessions("clover")
	username := session.GetString("username", "")
	approve := c.MustPostString("approve", "")
	resp := a.authServer.Authorize(c.Writer, c.Request, approve == "approve", username)
	resp.Write(c.Writer)
}

func (a *App) signinScreen(c *ace.C) {
	c.HTML("signin.html", nil)
}

func (a *App) signin(c *ace.C) {
	username := c.MustPostString("username", "")
	password := c.MustPostString("password", "")

	session := c.Sessions("clover")
	session.Set("username", username)
	session.Set("password", password)

	c.Redirect(c.MustQueryString("next", ""))
}

func (a *App) token(c *ace.C) {
	resp := a.authServer.Token(c.Writer, c.Request)
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
	_, resp := a.resourceServer.VerifyAccessToken(c.Writer, c.Request, "read_my_timeline")
	if resp != nil {
		resp.Write(c.Writer)
		return
	}

	c.HTML("home.html", nil)
}
