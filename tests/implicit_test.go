package clover

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildImplicitForm() url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("response_type", "token")

	return form
}

func TestImplicitAuthorize(t *testing.T) {
	store := newTestStore()
	c := newTestServer(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildImplicitForm())

	c.auth.Authorize(w, r, true).Write(w)

	token, err := getTokenFromUrl(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildVerifyForm(token))
	ac, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
	assert.Equal(t, ac.AccessToken, token)
}
