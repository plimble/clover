package clover

import (
	"github.com/plimble/clover"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func buildClientForm(responseType string, token string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("grant_type", "client_credentials")
	form.Set("client_secret", "xyz")
	form.Set("response_type", responseType)

	if token != "" {
		form.Set("access_token", token)
	}
	return form
}

func TestClientAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm("code", ""))
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.auth.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Authorize
	c.auth.Authorize(w, r, true)
	assert.Equal(t, 302, w.Code)

	// Get Token
	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientForm("token", ""))
	c.auth.Token(w, r)

	assert.Equal(t, 302, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)
	var resAt *clover.AccessToken

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientForm("token", token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}

func TestClientTokenAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm("token", ""))
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.auth.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Authorize
	c.auth.Authorize(w, r, true)

	assert.Equal(t, 302, w.Code)
	strings.Contains(w.HeaderMap["Location"][0], "access_token=")
	strings.Contains(w.HeaderMap["Location"][0], "expires=")

	token, err := getTokenFromUrl(w)
	assert.NoError(t, err)
	var resAt *clover.AccessToken

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientForm("token", token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}
