package tests

import (
	. "github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CodeSuite struct {
	suite.Suite
	authConfig     *AuthServerConfig
	resourceConfig *ResourceConfig
	authServer     *AuthServer
	resourceServer *ResourceServer
	store          *memory.Store
	client         *TestClient
}

func TestCodeSuite(t *testing.T) {
	suite.Run(t, &CodeSuite{})
}

func (t *CodeSuite) SetupSuite() {
	t.store = memory.New()

	t.authConfig = &AuthServerConfig{}
	t.resourceConfig = DefaultResourceConfig()

	t.authServer = NewAuthServer(t.store, t.authConfig)
	t.resourceServer = NewResourceServer(t.store, t.resourceConfig)

	t.authServer.AddGrantType(NewAuthorizationCode(t.store))
	t.authServer.AddRespType(NewCodeRespType(t.store, 5))
	t.authServer.SetAccessTokenRespType(NewAccessTokenRespType(t.store, t.store, 3600, 5000))

	t.client = &TestClient{
		ClientID:     "001",
		ClientSecret: "abc",
		GrantType:    []string{AUTHORIZATION_CODE},
		Scope:        []string{"read", "write"},
		RedirectURI:  "http://localhost/callback",
		Data: map[string]interface{}{
			"company_name": "xyz",
			"email":        "test@test.com",
		},
	}
}

func (t *CodeSuite) SetupTest() {
	t.store.SetClient(t.client)
}

func (t CodeSuite) TearDownTest() {
	t.store.Flush()
}

func (t *CodeSuite) defaultReqAuthorize() *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.ParseForm()
	r.PostForm.Set("response_type", "code")
	r.PostForm.Set("client_id", t.client.ClientID)

	return r
}

func (t *CodeSuite) defaultReqToken(code string) *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Set("Authorization", auth(t.client.ClientID, t.client.ClientSecret))
	r.ParseForm()
	r.PostForm.Set("grant_type", AUTHORIZATION_CODE)
	r.PostForm.Set("code", code)
	r.PostForm.Set("redirect_uri", t.client.RedirectURI)

	return r
}

func (t *CodeSuite) TestRequestAccessToken_Default() {
	//request authroize code
	r := t.defaultReqAuthorize()

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(existQuery(location, "code"))
	t.Equal("", w.Body.String())

	code := getQuery(location, "code")

	//reuest access token
	r = t.defaultReqToken(code)

	w = httptest.NewRecorder()
	resp = t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)
	t.Equal(3600, token.ExpiresIn)
	t.Len(token.Data, 0)
	t.NotEmpty(token.AccessToken)
	t.NotEmpty(token.RefreshToken)
	t.True(token.existScope("read"))
	t.True(token.existScope("write"))
}

func (t *CodeSuite) TestRequestAccessToken_CustomScope() {
	//request authroize code
	r := t.defaultReqAuthorize()
	r.PostForm.Set("scope", "read")

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(existQuery(location, "code"))
	t.Equal("", w.Body.String())

	code := getQuery(location, "code")

	//reuest access token
	r = t.defaultReqToken(code)

	w = httptest.NewRecorder()
	resp = t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)

	t.True(token.existScope("read"))
	t.False(token.existScope("write"))
}

func (t *CodeSuite) TestRequestAccessToken_RequiredState() {
	t.authConfig.StateParamRequired = true
	//request authroize code
	r := t.defaultReqAuthorize()
	r.PostForm.Set("scope", "read")
	r.PostForm.Set("state", "xxxxx")

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(existQuery(location, "code"))
	t.NotEmpty(existQuery(location, "state"))
	t.Equal("", w.Body.String())

	code := getQuery(location, "code")

	//reuest access token
	r = t.defaultReqToken(code)

	w = httptest.NewRecorder()
	resp = t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)

	t.True(token.existScope("read"))
	t.False(token.existScope("write"))
}

func (t *CodeSuite) TestVerifyAccessToken() {
	//request authroize code
	r := t.defaultReqAuthorize()
	t.authConfig.StateParamRequired = false

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Empty(w.Body.String())

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(existQuery(location, "code"))
	t.Equal("", w.Body.String())

	code := getQuery(location, "code")

	//reuest access token
	r = t.defaultReqToken(code)

	w = httptest.NewRecorder()
	resp = t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)
	t.Equal(3600, token.ExpiresIn)
	t.Len(token.Data, 0)
	t.NotEmpty(token.AccessToken)
	t.NotEmpty(token.RefreshToken)
	t.True(token.existScope("read"))
	t.True(token.existScope("write"))

	//verify access token
	r, _ = http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Set("access_token", token.AccessToken)
	r.URL.RawQuery = q.Encode()

	w = httptest.NewRecorder()
	at, resp := t.resourceServer.VerifyAccessToken(w, r, "read")
	t.Nil(resp)
	t.Equal(token.AccessToken, at.AccessToken)
}
