package verify

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/plimble/clover/oauth2"
	"github.com/plimble/clover/oauth2/mocks"
	"github.com/stretchr/testify/require"
)

func TestVerifyHandler(t *testing.T) {
	t.Run("ValidAccessToken", ValidAccessToken)
	t.Run("InValidAccessToken", InValidAccessToken)
}

func ValidAccessToken(t *testing.T) {
	storage := &mocks.Storage{}
	h := New(storage)
	res := httptest.NewRecorder()

	form := make(url.Values)
	form.Set("scope", "s1 s2")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "bearer a123456")

	at := &oauth2.AccessToken{
		AccessToken: "a123456",
		ClientID:    "c1",
		UserID:      "u1",
		Expired:     time.Now().Add(time.Second * 1000).Unix(),
		ExpiresIn:   1000,
		Scopes:      []string{"s1", "s2"},
		Extras: map[string]interface{}{
			"e1": "val1",
			"e2": "val1",
		},
	}

	storage.On("GetAccessToken", "a123456").Return(at, nil)

	h.ServeHTTP(res, req)

	require.Equal(t, 200, res.Code)
	require.Equal(t, "", res.Body.String())
	storage.AssertExpectations(t)
}

func InValidAccessToken(t *testing.T) {
	storage := &mocks.Storage{}
	h := New(storage)
	res := httptest.NewRecorder()

	form := make(url.Values)
	form.Set("scope", "s1 s2")
	req := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "bearer a123456")

	storage.On("GetAccessToken", "a123456").Return(nil, InvalidAccessToken("not found"))

	h.ServeHTTP(res, req)

	exp := map[string]interface{}{
		"error":             "invalid_accesstoken",
		"error_description": "not found",
	}

	resJson := make(map[string]interface{})
	err := json.NewDecoder(res.Body).Decode(&resJson)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, 400, res.Code)
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	require.Equal(t, exp, resJson)
	storage.AssertExpectations(t)
}
