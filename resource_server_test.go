package clover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockResourceServer struct {
	store *Mockallstore
}

func getTokenRequest() *http.Request {
	r, _ := http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Add("access_token", "123")
	r.URL.RawQuery = q.Encode()

	return r
}

func setupResourceServer() (*ResourceServer, *mockAuthCtrl) {
	store := NewMockallstore()
	config := NewResourceConfig(store)
	mock := &mockAuthCtrl{store}
	ctrl := NewResourceServer(config)

	return ctrl, mock
}

func TestResourceServer_GetTokenFromHttp_FromQuery(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Add("access_token", "123")
	r.URL.RawQuery = q.Encode()

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "123", token)
	assert.Nil(t, resp)
}

func TestResourceServer_GetTokenFromHttp_FromPost(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("POST", "/", nil)
	r.ParseForm()
	r.PostForm.Add("access_token", "123")

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "123", token)
	assert.Nil(t, resp)
}

func TestResourceServer_GetTokenFromHttp_FromAuthorization(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "Bearer 123")

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "123", token)
	assert.Nil(t, resp)
}

func TestResourceServer_GetTokenFromHttp_MultipleMethod(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "Bearer 123")
	r.ParseForm()
	r.PostForm.Add("access_token", "123")

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "", token)
	assert.Equal(t, errOnlyOneTokenMethod, resp)
}

func TestResourceServer_GetTokenFromHttp_InvalidMalFormed(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "XXX 123")

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "", token)
	assert.Equal(t, errMalFormedHeader, resp)
}

func TestResourceServer_GetTokenFromHttp_NoAccessToken(t *testing.T) {
	re, _ := setupResourceServer()

	r, _ := http.NewRequest("POST", "/", nil)

	token, resp := re.getTokenFromHttp(r)

	assert.Equal(t, "", token)
	assert.Equal(t, errNoTokenInRequest, resp)
}

func TestResourceServer_SetHeader(t *testing.T) {
	re, _ := setupResourceServer()

	testCases := []struct {
		resp   *Response
		scopes []string
		expWWW string
	}{
		{newResp(), nil, `Bearer realm="Service"`},
		{newResp(), []string{"1", "2"}, `Bearer realm="Service", scopes="1 2"`},
		{errParseURI, nil, `Bearer realm="Service", error="invalid_uri", error_description="Invalid parse uri"`},
		{errParseURI, []string{"1", "2"}, `Bearer realm="Service", scopes="1 2", error="invalid_uri", error_description="Invalid parse uri"`},
	}

	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		re.setHeader(testCase.resp, testCase.scopes, w)

		assert.Equal(t, w.Header().Get("WWW-Authenticate"), testCase.expWWW)
	}
}

func TestResourceServer_VerifyAccessToken_NoScope(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
	}

	m.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := re.VerifyAccessToken(w, r)

	assert.False(t, resp.IsError())
	assert.Equal(t, expat, at)
}

func TestResourceServer_VerifyAccessToken_TokenNotFound(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	m.store.On("GetAccessToken", "123").Return(nil, errors.New("not found"))

	at, resp := re.VerifyAccessToken(w, r)

	assert.Equal(t, errInvalidAccessToken, resp)
	assert.Nil(t, at)
}

func TestResourceServer_VerifyAccessToken_TokenExpire(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(0),
	}

	m.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := re.VerifyAccessToken(w, r)

	assert.Equal(t, errAccessTokenExpired, resp)
	assert.Nil(t, at)
}

func TestResourceServer_VerifyAccessToken_Scope(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2", "3"},
	}

	m.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := re.VerifyAccessToken(w, r, "1", "2")

	assert.False(t, resp.IsError())
	assert.Equal(t, expat, at)
}

func TestResourceServer_VerifyAccessToken_EmptyInToken(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
	}

	m.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := re.VerifyAccessToken(w, r, "1")

	assert.Equal(t, errInsufficientScope, resp)
	assert.Nil(t, at)
}

func TestResourceServer_VerifyAccessToken_NoScopeSupport(t *testing.T) {
	re, m := setupResourceServer()

	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2"},
	}

	m.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := re.VerifyAccessToken(w, r, "1", "3")

	assert.Equal(t, errInsufficientScope, resp)
	assert.Nil(t, at)
}
