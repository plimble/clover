package clover

import (
	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResourceServerSuite struct {
	suite.Suite
	store    *Mockallstore
	resource *ResourceServer
}

func TestResourceServer(t *testing.T) {
	suite.Run(t, &ResourceServerSuite{})
}

func (t *ResourceServerSuite) SetupTest() {
	t.store = NewMockallstore()
	config := &ResourceConfig{}
	t.resource = NewResourceServer(t.store, config)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_FromQuery() {
	r, _ := http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Add("access_token", "123")
	r.URL.RawQuery = q.Encode()

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("123", token)
	t.Nil(resp)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_FromPost() {
	r, _ := http.NewRequest("POST", "/", nil)
	r.ParseForm()
	r.PostForm.Add("access_token", "123")

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("123", token)
	t.Nil(resp)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_FromAuthorization() {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "Bearer 123")

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("123", token)
	t.Nil(resp)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_MultipleMethod() {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "Bearer 123")
	r.ParseForm()
	r.PostForm.Add("access_token", "123")

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("", token)
	t.Equal(errOnlyOneTokenMethod, resp)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_InvalidMalFormed() {
	r, _ := http.NewRequest("POST", "/", nil)
	r.Header.Add("Authorization", "XXX 123")

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("", token)
	t.Equal(errMalFormedHeader, resp)
}

func (t *ResourceServerSuite) TestGetTokenFromHttp_NoAccessToken() {
	r, _ := http.NewRequest("POST", "/", nil)

	token, resp := t.resource.getTokenFromHttp(r)

	t.Equal("", token)
	t.Equal(errNoTokenInRequest, resp)
}

func (t *ResourceServerSuite) TestSetHeader() {
	testCases := []struct {
		resp   *Response
		scopes []string
		expWWW string
	}{
		{&Response{code: 200}, nil, `Bearer realm="Service"`},
		{&Response{code: 200}, []string{"1", "2"}, `Bearer realm="Service", scopes="1 2"`},
		{errParseURI, nil, `Bearer realm="Service", error="invalid_uri", error_description="Invalid parse uri"`},
		{errParseURI, []string{"1", "2"}, `Bearer realm="Service", scopes="1 2", error="invalid_uri", error_description="Invalid parse uri"`},
	}

	for _, testCase := range testCases {
		w := httptest.NewRecorder()
		t.resource.setHeader(testCase.resp, testCase.scopes, w)

		t.Equal(w.Header().Get("WWW-Authenticate"), testCase.expWWW)
	}
}

func (t *ResourceServerSuite) TestVerifyAccessToken_NoScope() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
	}

	t.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := t.resource.VerifyAccessToken(w, r)

	t.False(resp.IsError())
	t.Equal(expat, at)
}

func (t *ResourceServerSuite) TestVerifyAccessToken_TokenNotFound() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	t.store.On("GetAccessToken", "123").Return(nil, errors2.NewAnyError())

	at, resp := t.resource.VerifyAccessToken(w, r)

	t.Equal(errInvalidAccessToken, resp)
	t.Nil(at)
}

func (t *ResourceServerSuite) TestVerifyAccessToken_TokenExpire() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(-1),
	}

	t.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := t.resource.VerifyAccessToken(w, r)

	t.Equal(errAccessTokenExpired, resp)
	t.Nil(at)
}

func (t *ResourceServerSuite) TestVerifyAccessToken_Scope() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2", "3"},
	}

	t.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := t.resource.VerifyAccessToken(w, r, "1", "2")

	t.False(resp.IsError())
	t.Equal(expat, at)
}

func (t *ResourceServerSuite) TestVerifyAccessToken_EmptyInToken() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
	}

	t.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := t.resource.VerifyAccessToken(w, r, "1")

	t.Equal(errInsufficientScope, resp)
	t.Nil(at)
}

func (t *ResourceServerSuite) TestVerifyAccessToken_NoScopeSupport() {
	w := httptest.NewRecorder()
	r := getTokenRequest()

	expat := &AccessToken{
		AccessToken: "123",
		Expires:     addSecondUnix(60),
		Scope:       []string{"1", "2"},
	}

	t.store.On("GetAccessToken", "123").Return(expat, nil)

	at, resp := t.resource.VerifyAccessToken(w, r, "1", "3")

	t.Equal(errInsufficientScope, resp)
	t.Nil(at)
}

func getTokenRequest() *http.Request {
	r, _ := http.NewRequest("GET", "/", nil)
	q := r.URL.Query()
	q.Add("access_token", "123")
	r.URL.RawQuery = q.Encode()

	return r
}
