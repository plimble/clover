package main

import (
	"github.com/gorilla/sessions"
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/pongo2"
	"github.com/plimble/aero"
	"github.com/plimble/clover"
	aeroClover "github.com/plimble/clover/store/aerospike"
	"gopkg.in/mgo.v2"
	"net/url"
)

func setupStore(client *aero.Client) *aeroClover.AeroStore {
	aeroClover.New(asClient, ns, tokenLifeTime, authCodeLifetime, refresLifeTime)
	return m
}

func getUser(session *mgo.Session, db string) mongo.GetUserFunc {
	return func(username, password string) (string, []string, error) {
		return "1", nil, nil
	}

}

func getClient(session *mgo.Session, db string) mongo.GetClientFunc {
	return func(clientID string) (clover.Client, error) {
		var c clover.DefaultClient
		if err := session.DB(db).C("oauth_client").FindId(clientID).One(&c); err != nil {
			return nil, err
		}

		return &c, nil
	}
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
	session, err := mgo.Dial("172.17.8.101:27017")
	if err != nil {
		panic(err)
	}

	store := setupStore(session)
	authConfig := clover.NewAuthConfig(store)
	authConfig.AllowImplicit = true
	authConfig.AddClientGrant()
	authConfig.AddPasswordGrant(store)
	authConfig.AddRefreshGrant(store)
	authConfig.AddAuthCodeGrant(store)
	authConfig.SetDefaultScopes("read_my_timeline", "read_my_friend")
	auth := clover.NewAuthServer(authConfig)

	resourceconfig := clover.NewResourceConfig(store)
	resource := clover.NewResourceServer(resourceconfig)

	app := &App{
		auth:     auth,
		resource: resource,
	}

	a := ace.New()

	a.HtmlTemplate(pongo2.Pongo2(&pongo2.TemplateOptions{
		Directory: "views",
	}))

	cookie := sessions.NewCookieStore([]byte("secret"))
	a.Session("clover", cookie, nil)

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
