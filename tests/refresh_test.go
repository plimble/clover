package clover

import (
	"encoding/json"
	"github.com/plimble/clover"
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
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildPasswordForm())

	// Get Token
	c.auth.Token(w, r)
	assert.Equal(t, 200, w.Code)
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resJSON)
	assert.NoError(t, err)

	// Get Refresh Token
	w = httptest.NewRecorder()
	r = newTestRequest("http://localhost", "", buildRefreshTokenForm(resJSON["refresh_token"].(string)))
	c.auth.Token(w, r)

	var resAt *clover.AccessToken
	token, err := getTokenFromBody(w)

	r = newTestRequest("http://localhost", "", buildClientTokenForm(token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})

	assert.NotNil(t, resAt)
}
