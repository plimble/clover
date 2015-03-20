package clover

import (
	"github.com/plimble/clover"
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
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildPasswordForm())
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.auth.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Get Token
	c.auth.Token(w, r)

	assert.Equal(t, 200, w.Code)
	validateResponseToken(t, w.Body.String())

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)
	var resAt *clover.AccessToken

	r = newTestRequest("http://localhost", "", buildClientTokenForm(token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}
