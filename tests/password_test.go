package clover

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildPasswordForm() url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("grant_type", "password")
	form.Set("client_secret", "xyz")
	form.Set("response_type", "token")
	form.Set("username", "test")
	form.Set("password", "1234")

	return form
}

func TestPasswordAuthorize(t *testing.T) {
	store := newTestStore()
	c := newTestServer(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildPasswordForm())

	// Get Token
	c.auth.Token(w, r).Write(w)

	assert.Equal(t, 200, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildAuthTokenForm(token))
	ac, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, ac.AccessToken, token)
}
