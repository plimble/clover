package tests

import (
	. "github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JWTAccessTokenSuite struct {
	suite.Suite
	authConfig     *AuthServerConfig
	resourceConfig *ResourceConfig
	authServer     *AuthServer
	resourceServer *ResourceServer
	store          *memory.Store
	client         *TestClient
}

func TestJWTAccessTokenSuite(t *testing.T) {
	suite.Run(t, &JWTAccessTokenSuite{})
}

func (t *JWTAccessTokenSuite) SetupSuite() {
	t.store = memory.New()

	t.authConfig = &AuthServerConfig{}
	t.resourceConfig = DefaultResourceConfig()

	t.authServer = NewAuthServer(t.store, t.authConfig)
	t.resourceServer = NewResourceServer(t.store, t.resourceConfig)
	t.resourceServer.UseJWTAccessToken(t.store)

	t.authServer.AddGrantType(NewClientCredential(t.store))
	t.authServer.SetAccessTokenRespType(NewJWTAccessTokenRespType(t.store, t.store, 3600, 5000))

	t.client = &TestClient{
		ClientID:     "001",
		ClientSecret: "abc",
		GrantType:    []string{CLIENT_CREDENTIALS},
		Scope:        []string{"read", "write"},
		RedirectURI:  "http://localhost/callback",
		Data: map[string]interface{}{
			"company_name": "xyz",
			"email":        "test@test.com",
		},
	}
}

func (t *JWTAccessTokenSuite) SetupTest() {
	t.store.SetClient(t.client)
	t.store.AddHSKey(t.client.ClientID)
}

func (t JWTAccessTokenSuite) TearDownTest() {
	t.store.Flush()
}

func (t *JWTAccessTokenSuite) defaultReqToken() *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Set("Authorization", auth(t.client.ClientID, t.client.ClientSecret))
	r.ParseForm()
	r.PostForm.Set("grant_type", CLIENT_CREDENTIALS)

	return r
}

func (t *JWTAccessTokenSuite) TestRequestAccessToken_Default() {
	//reuest access token
	r := t.defaultReqToken()

	w := httptest.NewRecorder()
	resp := t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)
	t.Equal(3600, token.ExpiresIn)
	t.Len(token.Data, 2)
	t.NotEmpty(token.AccessToken)
	t.True(token.existScope("read"))
	t.True(token.existScope("write"))
}

func (t *JWTAccessTokenSuite) TestRequestAccessToken_CustomScope() {
	//reuest access token
	r := t.defaultReqToken()
	r.PostForm.Set("scope", "read")

	w := httptest.NewRecorder()
	resp := t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)

	t.True(token.existScope("read"))
	t.False(token.existScope("write"))
}

func (t *JWTAccessTokenSuite) TestVerifyAccessToken() {
	//reuest access token
	r := t.defaultReqToken()
	r.PostForm.Set("scope", "read")

	w := httptest.NewRecorder()
	resp := t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)

	t.True(token.existScope("read"))
	t.False(token.existScope("write"))

	//verify access token
	r, _ = http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Set("access_token", token.AccessToken)
	r.URL.RawQuery = q.Encode()

	w = httptest.NewRecorder()
	at, resp := t.resourceServer.VerifyAccessToken(w, r, "read")
	t.False(resp.IsError())
	t.NotNil(at)
}
