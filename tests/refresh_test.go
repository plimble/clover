package clover

import (
	"encoding/json"
	"github.com/plimble/clover"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildRefreshForm(responseType, grantType, refreshToken, user, pass string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("grant_type", grantType)
	form.Set("client_secret", "xyz")
	form.Set("response_type", responseType)
	form.Set("refresh_token", refreshToken)
	form.Set("username", user)
	form.Set("password", pass)

	return form
}

func TestRefreshToken(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildRefreshForm("code", "password", "", "test", "1234"))
	fn := func(client clover.Client, scopes []string) {}

	// Validate Authorize
	c.auth.ValidateAuthorize(w, r, fn)
	assert.Equal(t, 200, w.Code)

	// Get Token
	c.auth.Token(w, r)
	assert.Equal(t, 200, w.Code)
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resJSON)
	assert.NoError(t, err)

	// Get Refresh Token
	w = httptest.NewRecorder()
	r = newTestRequest("http://localhost", "", buildRefreshForm("token", "refresh_token", resJSON["refresh_token"].(string), "test", "1234"))
	c.auth.Token(w, r)

	var resAt *clover.AccessToken
	token, err := getTokenFromBody(w)

	r = newTestRequest("http://localhost", "", buildClientForm("token", token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}
