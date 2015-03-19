package clover

import (
	"bytes"
	"encoding/json"
	"github.com/plimble/clover"
	"github.com/plimble/clover/store/memory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func newTestServer() *clover.AuthorizeServer {
	config := &clover.Config{
		AccessLifeTime:       1,
		AuthCodeLifetime:     60,
		RefreshTokenLifetime: 1,
		AllowCredentialsBody: true,
		AllowImplicit:        true,
		StateParamRequired:   false,
	}

	auth := clover.NewAuthServer(newTestStore(), config)
	auth.RegisterClientGrant()
	auth.RegisterPasswordGrant()
	auth.RegisterRefreshGrant()
	auth.RegisterAuthCodeGrant()
	auth.RegisterImplicitGrant()
	auth.SetDefaultScopes("read_my_timeline", "read_my_friend")
	return auth
}

func newTestStore() clover.AuthServerStore {
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

func validateResponseToken(t *testing.T, body string) {
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(body), &resJSON)
	assert.NoError(t, err)
	assert.Equal(t, resJSON["token_type"], "bearer")
	assert.Equal(t, resJSON["expires_in"], 1)
	assert.NotEmpty(t, resJSON["access_token"])
	// assert.NotEmpty(t, resJSON["refresh_token"])
}
