package clover

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildRefreshTokenForm(refreshToken string) url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("grant_type", "refresh_token")
	form.Set("client_secret", "xyz")
	form.Set("refresh_token", refreshToken)

	return form
}

func TestRefreshToken(t *testing.T) {
	store := newTestStore()
	c := newTestServer(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildPasswordForm())

	// Get Token
	c.auth.Token(w, r).Write(w)
	assert.Equal(t, 200, w.Code)

	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resJSON)
	assert.NoError(t, err)

	// Get Refresh Token
	w = httptest.NewRecorder()
	r = newTestRequest("http://localhost", "", buildRefreshTokenForm(resJSON["refresh_token"].(string)))
	c.auth.Token(w, r).Write(w)
	assert.Equal(t, 200, w.Code)

	token, err := getTokenFromBody(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildAuthTokenForm(token))
	ac, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, ac.AccessToken, token)
}
