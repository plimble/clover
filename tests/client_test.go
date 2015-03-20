package clover

import (
	"github.com/plimble/clover"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func buildClientForm() url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("response_type", "code")

	return form
}

func TestClientAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientForm())
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.auth.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Authorize
	c.auth.Authorize(w, r, true)
	assert.Equal(t, 302, w.Code)

	// Get Token
	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientTokenForm(""))
	c.auth.Token(w, r)

	assert.Equal(t, 302, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)
	var resAt *clover.AccessToken

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientTokenForm(token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}

func TestClientTokenAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildClientTokenForm(""))
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

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientTokenForm(token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}
