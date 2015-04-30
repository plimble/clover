package tests

import (
	. "github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ImplicitSuite struct {
	suite.Suite
	authConfig     *AuthServerConfig
	resourceConfig *ResourceConfig
	authServer     *AuthServer
	resourceServer *ResourceServer
	store          *memory.Store
	client         *TestClient
}

func TestImplicitSuite(t *testing.T) {
	suite.Run(t, &ImplicitSuite{})
}

func (t *ImplicitSuite) SetupSuite() {
	t.store = memory.New()

	t.authConfig = &AuthServerConfig{}
	t.resourceConfig = DefaultResourceConfig()

	t.authServer = NewAuthServer(t.store, t.store, t.authConfig)
	t.resourceServer = NewResourceServer(t.store, t.resourceConfig)

	t.authServer.AddGrantType(NewAuthorizationCode(t.store))
	t.authServer.AddRespType(NewImplicitRespType(t.store, t.store, 3600, 5000))

	t.client = &TestClient{
		ClientID:     "001",
		ClientSecret: "abc",
		GrantType:    []string{IMPLICIT},
		Scope:        []string{"read", "write"},
		RedirectURI:  "http://localhost/callback",
		Data: map[string]interface{}{
			"company_name": "xyz",
			"email":        "test@test.com",
		},
	}
}

func (t *ImplicitSuite) SetupTest() {
	t.store.SetClient(t.client)
}

func (t ImplicitSuite) TearDownTest() {
	t.store.Flush()
}

func (t *ImplicitSuite) defaultReqAuthorize() *http.Request {
	r, _ := http.NewRequest("POST", "/", nil)
	r.ParseForm()
	r.PostForm.Set("response_type", "token")
	r.PostForm.Set("client_id", t.client.ClientID)

	return r
}

func (t *ImplicitSuite) TestRequestAccessToken_Default() {
	//request authroize code
	r := t.defaultReqAuthorize()

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(getFragment(location))
	t.Equal("", w.Body.String())
}

func (t *ImplicitSuite) TestRequestAccessToken_CustomScope() {
	//request authroize code
	r := t.defaultReqAuthorize()
	r.PostForm.Set("scope", "read")

	w := httptest.NewRecorder()

	resp := t.authServer.Authorize(w, r, true, "userid")
	resp.Write(w)

	location := w.Header().Get("LOCATION")

	t.Equal(302, w.Code)
	t.Equal(t.client.RedirectURI, getURL(location))
	t.NotEmpty(getFragment(location))
	t.Equal("", w.Body.String())
}

func (t *ImplicitSuite) TestRequestAccessToken_RequiredState() {
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
	t.NotEmpty(getFragment(location))
	t.Equal("", w.Body.String())
}
