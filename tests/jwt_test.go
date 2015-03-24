package clover

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func buildJWTForm() url.Values {
	form := url.Values{}
	form.Set("client_id", "1001")
	form.Set("client_secret", "xyz")
	form.Set("grant_type", "client_credentials")
	form.Set("username", "test")
	form.Set("password", "1234")

	return form
}

func getTokenJWTFromBody(w *httptest.ResponseRecorder) (string, error) {
	var resJSON map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resJSON)
	return resJSON["access_token"].(string), err
}

func TestJWT(t *testing.T) {
	store := newTestStore()
	store.AddHSKey("1001")

	c := newTestServerJWT(store)

	w := httptest.NewRecorder()
	r := newTestRequest("http://localhost", "", buildJWTForm())

	c.auth.Token(w, r).Write(w)

	token, err := getTokenJWTFromBody(w)
	assert.NoError(t, err)

	r = newTestRequest("http://localhost", "", buildVerifyForm(token))
	_, resp := c.resource.VerifyAccessToken(w, r, "read_my_timeline")
	assert.False(t, resp.IsError())
	assert.False(t, resp.IsRedirect())
}
