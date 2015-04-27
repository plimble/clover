package tests

import (
	. "github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ClientSuite struct {
	suite.Suite
	authConfig     *AuthServerConfig
	resourceConfig *ResourceConfig
	authServer     *AuthServer
	resourceServer *ResourceServer
	store          *memory.Store
	client         *TestClient
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, &ClientSuite{})
}

func (t *ClientSuite) SetupSuite() {
	t.store = memory.New()

	t.authConfig = &AuthServerConfig{}
	t.resourceConfig = &ResourceConfig{}

	t.authServer = NewAuthServer(t.store, t.authConfig)
	t.resourceServer = NewResourceServer(t.store, t.resourceConfig)

	t.authServer.AddGrantType(NewClientCredential(t.store))
	t.authServer.SetAccessTokenRespType(NewAccessTokenRespType(t.store, t.store, 3600, 5000))

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

func (t *ClientSuite) SetupTest() {
	t.store.SetClient(t.client)
}

func (t ClientSuite) TearDownTest() {
	t.store.Flush()
}

func (t *ClientSuite) defaultReqToken() *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Set("Authorization", auth(t.client.ClientID, t.client.ClientSecret))
	r.ParseForm()
	r.PostForm.Set("grant_type", CLIENT_CREDENTIALS)

	return r
}

func (t *ClientSuite) TestRequestAccessToken_Default() {
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

func (t *ClientSuite) TestRequestAccessToken_CustomScope() {
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
