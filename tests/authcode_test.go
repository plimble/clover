package clover

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildAuthCodeForm(responseType string) url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("response_type", responseType)

	return form
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

func TestCodeAuthorize(t *testing.T) {
	s := newTestStore()
	c := newTestServer(s)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildAuthCodeForm("code"))

	// Authorize
	c.auth.Authorize(w, r, true).Write(w)
	assert.Equal(t, 302, w.Code)

	// Get Token
	auth_code, err := getTokenFromUrl(w)
	assert.NoError(t, err)

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildAuthTokenForm(auth_code))
	c.auth.Token(w, r).Write(w)

	assert.Equal(t, 302, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildAuthTokenForm(token))
	ac, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, ac.AccessToken, token)
}
