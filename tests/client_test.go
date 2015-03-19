package clover

import (
	"github.com/plimble/clover"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func buildClientForm(responseType string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("grant_type", "client_credentials")
	form.Set("client_secret", "xyz")
	form.Set("response_type", responseType)

	return form
}

func TestClientAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm("code"))
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Authorize
	c.Authorize(w, r, true)
	assert.Equal(t, 302, w.Code)

	// Get Token
	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientForm("token"))
	c.Token(w, r)

	assert.Equal(t, 302, w.Code)
	validateResponseToken(t, w.Body.String())
}

func TestClientTokenAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm("token"))
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Authorize
	c.Authorize(w, r, true)

	assert.Equal(t, 302, w.Code)
	strings.Contains(w.HeaderMap["Location"][0], "access_token=")
	strings.Contains(w.HeaderMap["Location"][0], "expires=")
}
