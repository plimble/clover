package main

import (
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/pongo2"
	"github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/plimble/clover/tests"
	"github.com/plimble/sessions/store/cookie"
	"net/url"
)

func isLogin(c *ace.C) {
	session := c.Sessions("clover")
	if session.IsNew {
		c.Redirect("/signin?next=" + url.QueryEscape(c.Request.RequestURI))
		c.Abort()
		return
	}
	c.Next()
}

func initData(store *memory.Store) {
	client := &tests.TestClient{
		ClientID:     "001",
		ClientSecret: "abc",
		GrantType:    []string{clover.AUTHORIZATION_CODE, clover.IMPLICIT, clover.CLIENT_CREDENTIALS, clover.PASSWORD, clover.REFRESH_TOKEN},
		Scope:        []string{"read", "write"},
		RedirectURI:  "http://localhost:4000/callback",
		Data: map[string]interface{}{
			"company_name": "xyz",
			"email":        "test@test.com",
		},
	}

	user := &tests.TestUser{
		ID:       "111",
		Username: "test",
		Password: "1234",
		Data: map[string]interface{}{
			"email": "test@test.com",
		},
	}

	store.SetClient(client)
	store.SetUser(user)
}

func main() {
	store := memory.New()
	initData(store)

	authServer := clover.NewAuthServer(store, clover.DefaultAuthServerConfig())
	authServer.AddGrantType(clover.NewClientCredential(store))
	authServer.AddGrantType(clover.NewPassword(store))
	authServer.AddGrantType(clover.NewAuthorizationCode(store))
	authServer.AddGrantType(clover.NewRefreshToken(store))
	authServer.AddRespType(clover.NewCodeRespType(store, 500))
	authServer.AddRespType(clover.NewImplicitRespType(store, store, 3600, 5000))
	authServer.SetAccessTokenRespType(clover.NewAccessTokenRespType(store, store, 3600, 5000))

	resourceServer := clover.NewResourceServer(store, clover.DefaultResourceConfig())

	app := &App{
		authServer:     authServer,
		resourceServer: resourceServer,
	}

	a := ace.New()

	a.HtmlTemplate(pongo2.Pongo2(&pongo2.TemplateOptions{
		Directory: "views",
	}))

	a.Use(ace.Session(cookie.NewCookieStore(), nil))

	//http://localhost:4000/oauth?client_id=1001&response_type=code
	a.GET("/oauth", isLogin, app.grantScreen)
	a.POST("/oauth", isLogin, app.grant)
	a.GET("/signin", app.signinScreen)
	a.POST("/signin", app.signin)
	a.POST("/token", app.token)
	a.GET("/callback", app.callback)
	a.GET("/home", isLogin, app.home)

	a.Run(":4000")
}
