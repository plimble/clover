package clover

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func buildClientForm() url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("response_type", "code")

	return form
}

func TestClientAuthorize(t *testing.T) {
	store := newTestStore()
	c := newTestServer(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm())

	// Get Token
	r = newTestRequest("http://localhost", "", buildClientTokenForm(""))
	c.auth.Token(w, r).Write(w)

	assert.Equal(t, 200, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildVerifyForm(token))
	at, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, at.AccessToken, token)
}

func TestClientTokenAuthorize(t *testing.T) {
	store := newTestStore()
	c := newTestServer(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientTokenForm(""))

	// Authorize
	c.auth.Authorize(w, r, true).Write(w)

	assert.Equal(t, 302, w.Code)
	strings.Contains("http://localhost", "access_token=")
	strings.Contains("http://localhost", "expires=")

	token, err := getTokenFromUrl(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildVerifyForm(token))
	ac, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, ac.AccessToken, token)
}
