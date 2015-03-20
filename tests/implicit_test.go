package clover

import (
	"github.com/plimble/clover"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildImplicitForm(responseType, token string) url.Values {
	form := url.Values{}
	form.Set("redirect_uri", "http://localhost:4000/callback")
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("response_type", responseType)

	if token != "" {
		form.Set("access_token", token)
	}
	return form
}

func TestImplicitAuthorize(t *testing.T) {
	c := newTestServer()

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildImplicitForm("token", ""))

	c.auth.Authorize(w, r, true)

	token, err := getTokenFromUrl(w)
	assert.NoError(t, err)
	var resAt *clover.AccessToken

	r = newTestRequest(w.HeaderMap["Location"][0], "", buildClientForm("token", token))
	c.resource.VerifyAccessToken(w, r, []string{"read_my_timeline"}, func(at *clover.AccessToken) {
		resAt = at
	})
	assert.NotNil(t, resAt)
}
