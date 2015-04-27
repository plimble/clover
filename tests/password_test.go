package tests

import (
	. "github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PasswordSuite struct {
	suite.Suite
	authConfig     *AuthServerConfig
	resourceConfig *ResourceConfig
	authServer     *AuthServer
	resourceServer *ResourceServer
	store          *memory.Store
	client         *TestClient
	user           *TestUser
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, &PasswordSuite{})
}

func (t *PasswordSuite) SetupSuite() {
	t.store = memory.New()

	t.authConfig = &AuthServerConfig{}
	t.resourceConfig = DefaultResourceConfig()

	t.authServer = NewAuthServer(t.store, t.authConfig)
	t.resourceServer = NewResourceServer(t.store, t.resourceConfig)

	t.authServer.AddGrantType(NewPassword(t.store))
	t.authServer.SetAccessTokenRespType(NewAccessTokenRespType(t.store, t.store, 3600, 5000))

	t.client = &TestClient{
		ClientID:     "001",
		ClientSecret: "abc",
		GrantType:    []string{PASSWORD},
		Scope:        []string{"read", "write"},
		RedirectURI:  "http://localhost/callback",
		Data: map[string]interface{}{
			"company_name": "xyz",
			"email":        "test@test.com",
		},
	}

	t.user = &TestUser{
		ID:       "111",
		Username: "tester",
		Password: "1234",
		Data: map[string]interface{}{
			"email": "test@test.com",
		},
	}
}

func (t *PasswordSuite) SetupTest() {
	t.store.SetClient(t.client)
	t.store.SetUser(t.user)
}

func (t PasswordSuite) TearDownTest() {
	t.store.Flush()
}

func (t *PasswordSuite) defaultReqToken() *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Set("Authorization", auth(t.client.ClientID, t.client.ClientSecret))
	r.ParseForm()
	r.PostForm.Set("username", t.user.Username)
	r.PostForm.Set("password", t.user.Password)
	r.PostForm.Set("grant_type", PASSWORD)

	return r
}

func (t *PasswordSuite) TestRequestAccessToken_Default() {
	//reuest access token
	r := t.defaultReqToken()

	w := httptest.NewRecorder()
	resp := t.authServer.Token(w, r)
	resp.Write(w)

	token := getToken(w.Body.Bytes())

	t.Equal(200, w.Code)
	t.Equal(3600, token.ExpiresIn)
	t.Len(token.Data, 1)
	t.NotEmpty(token.AccessToken)
	t.True(token.existScope("read"))
	t.True(token.existScope("write"))
}

func (t *PasswordSuite) TestRequestAccessToken_CustomScope() {
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

func (t *PasswordSuite) TestVerifyAccessToken() {
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
	t.Nil(resp)
	t.Equal(token.AccessToken, at.AccessToken)
}
