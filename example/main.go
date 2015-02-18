package main

import (
	"github.com/gorilla/sessions"
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/pongo2"
	"github.com/plimble/clover"
	"github.com/plimble/clover/store/mongo"
	"gopkg.in/mgo.v2"
	"net/url"
)

func setUpStore() clover.Store {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	session.DB("oauth").C("oauth_client").UpsertId("1001", clover.Client{
		ClientID:     "1001",
		ClientSecret: "xyz",
		GrantType:    []string{clover.AUTHORIZATION_CODE, clover.PASSWORD, clover.CLIENT_CREDENTIALS, clover.REFRESH_TOKEN},
		UserID:       "1",
		Scope:        []string{"read_my_timeline", "read_my_friend"},
		RedirectURI:  "http://localhost:4000/callback",
	})

	session.DB("oauth").C("oauth_scope").UpsertId("post_my_wall", clover.Scope{"post_my_wall", "Can post my wall"})
	session.DB("oauth").C("oauth_scope").UpsertId("read_my_timeline", clover.Scope{"read_my_timeline", "Can read my timeline"})
	session.DB("oauth").C("oauth_scope").UpsertId("read_my_friend", clover.Scope{"read_my_friend", "Can read my friend"})

	return mongo.New(session, "oauth")
}

func isLogin(c *ace.C) {
	if c.Session.IsNew() {
		c.Redirect("/signin?next=" + url.QueryEscape(c.Request.RequestURI))
		c.Abort()
		return
	}
	c.Next()
}

func main() {

	config := clover.DefaultConfig()
	config.Store = setUpStore()
	config.Grants = []string{clover.CLIENT_CREDENTIALS, clover.PASSWORD, clover.AUTHORIZATION_CODE, clover.REFRESH_TOKEN}
	config.AllowImplicit = true
	cv := clover.New(config)
	cv.SetDefaultScopes("read_my_timeline", "read_my_friend")

	a := ace.New()
	a.UseHtmlTemplate(pongo2.Pongo2(&pongo2.TemplateOptions{
		Directory: "views",
	}))

	cookie := sessions.NewCookieStore([]byte("secret"))
	a.UseSession("clover", cookie, nil)

	a.GET("/oauth", isLogin, func(c *ace.C) {
		ar := cv.ValidateAuthorize(c.Writer, c.Request)
		//if validate failed
		if ar == nil {
			return
		}

		//go to dialog page
		descScopes, err := cv.GetScopeDescription(ar.Scope)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.HTML("authorize.html", map[string]interface{}{
			"auth":       ar,
			"descScopes": descScopes,
		})
	})

	a.POST("/oauth", isLogin, func(c *ace.C) {
		approve := c.Request.FormValue("approve")
		cv.Authorize(c.Writer, c.Request, approve == "approve")
	})

	a.GET("/signin", func(c *ace.C) {
		c.HTML("signin.html", nil)
	})

	a.POST("/signin", func(c *ace.C) {
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")
		c.Session.SetString("username", username)
		c.Session.SetString("password", password)

		c.Redirect(c.Request.FormValue("next"))
	})

	a.POST("/token", func(c *ace.C) {
		cv.Token(c.Writer, c.Request)
	})

	a.GET("/callback", func(c *ace.C) {
		if c.Request.FormValue("err") != "" {
			c.String(500, "%s %s", c.Request.FormValue("err"), c.Request.FormValue("desc"))
		} else {
			c.String(200, c.Request.FormValue("code"))
		}
	})

	a.Run(":4000")
}
