package clover

import (
	"bytes"
	"encoding/json"
	"github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

type testApp struct {
	auth     *clover.AuthServer
	resource *clover.ResourceServer
}

func newTestServer(store *memory.Storage) *testApp {
	config := clover.NewAuthConfig(store)
	config.AddClientGrant()
	config.AddPasswordGrant(store)
	config.AddRefreshGrant(store)
	config.AddAuthCodeGrant(store)

	config.AllowCredentialsBody = true
	config.AllowImplicit = true
	config.AccessLifeTime = 1
	config.AuthCodeLifetime = 60
	config.RefreshTokenLifetime = 30

	config.SetDefaultScopes("read_my_timeline", "read_my_friend")
	auth := clover.NewAuthServer(config)

	resourceconfig := clover.NewResourceConfig(store)
	resource := clover.NewResourceServer(resourceconfig)

	return &testApp{auth, resource}
}

func newTestServerJWT(store *memory.Storage) *testApp {
	config := clover.NewAuthConfig(store)
	config.AddClientGrant()
	config.AddPasswordGrant(store)
	config.AddRefreshGrant(store)
	config.AddAuthCodeGrant(store)
	config.UseJWTAccessTokens(store)

	config.AllowCredentialsBody = true
	config.AllowImplicit = true
	config.AccessLifeTime = 1
	config.AuthCodeLifetime = 60
	config.RefreshTokenLifetime = 30

	config.SetDefaultScopes("read_my_timeline", "read_my_friend")
	auth := clover.NewAuthServer(config)

	resourceconfig := clover.NewResourceConfig(store)
	resourceconfig.UseJWTAccessTokens(store)

	resource := clover.NewResourceServer(resourceconfig)

	return &testApp{auth, resource}
}

func newTestStore() *memory.Storage {
	// New Store(Memory)
	store := memory.New()

	// Add User
	store.User = make(map[string]*memory.User)
	store.User["test"] = &memory.User{"test", "1234"}

	// Set Config
	sc := make(map[string]*clover.DefaultClient)
	sc["1001"] = &clover.DefaultClient{
		ClientID:     "1001",
		ClientSecret: "xyz",
		GrantType:    []string{clover.AUTHORIZATION_CODE, clover.PASSWORD, clover.CLIENT_CREDENTIALS, clover.REFRESH_TOKEN, clover.IMPLICIT},
		UserID:       "1",
		Scope:        []string{"read_my_timeline", "read_my_friend"},
		RedirectURI:  "http://localhost:4000/callback",
	}

	store.Client = sc
	return store
}

func newTestRequest(urlRequest, authType string, form url.Values) *http.Request {
	r, err := http.NewRequest("POST", urlRequest, bytes.NewBufferString(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if authType != "" {
		r.Header.Add("Authorization", authType)
	}

	if err != nil {
		panic(err)
	}
	return r
}

func getTokenFromUrl(w *httptest.ResponseRecorder) (string, error) {
	str := strings.Split(w.HeaderMap["Location"][0], "=")
	str = strings.Split(str[1], "&")
	return str[0], nil
}

func getTokenFromBody(w *httptest.ResponseRecorder) (string, error) {
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resJSON)
	return resJSON["access_token"].(string), err
}

func buildAuthTokenForm(token string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("grant_type", "authorization_code")

	if token != "" {
		form.Set("access_token", token)
	}
	return form
}

func buildVerifyForm(token string) url.Values {
	form := url.Values{}
	form.Set("access_token", token)
	return form
}

func buildClientTokenForm(token string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("grant_type", "client_credentials")
	form.Set("client_secret", "xyz")
	form.Set("response_type", "token")

	if token != "" {
		form.Set("access_token", token)
	}
	return form
}

func validateResponseToken(t *testing.T, body string) {
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(body), &resJSON)
	assert.NoError(t, err)
	assert.Equal(t, resJSON["token_type"], "bearer")
	assert.Equal(t, resJSON["expires_in"], 1)
	assert.NotEmpty(t, resJSON["access_token"])
}
