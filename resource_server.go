package clover

import (
	"github.com/plimble/utils/strings2"
	"net/http"
	"strings"
)

type ResourceConfig struct {
	WWWRealm              string
	TokenBearerHeaderName string
	TokenParamName        string
}

func DefaultResourceConfig() *ResourceConfig {
	return &ResourceConfig{
		WWWRealm:              "Service",
		TokenBearerHeaderName: "Bearer",
		TokenParamName:        "access_token",
	}
}

type ResourceServer struct {
	config           *ResourceConfig
	accessTokenStore AccessTokenStore
}

func NewResourceServer(accessTokenStore AccessTokenStore, config *ResourceConfig) *ResourceServer {
	if config.WWWRealm == "" {
		config.WWWRealm = "Service"
	}

	return &ResourceServer{config, accessTokenStore}
}

func (s *ResourceServer) UseJWTAccessToken(publicKeyStore PublicKeyStore) {
	s.accessTokenStore = newJWTTokenStore(publicKeyStore)
}

func (s *ResourceServer) VerifyAccessToken(w http.ResponseWriter, r *http.Request, scopes ...string) (*AccessToken, *Response) {
	token, resp := s.getTokenFromHttp(r)
	if resp != nil {
		return nil, s.setHeader(resp, nil, w)
	}

	at, err := s.accessTokenStore.GetAccessToken(token)
	if err != nil {
		return nil, s.setHeader(errInvalidAccessToken, nil, w)
	}

	if at.Expires != 0 && isExpireUnix(at.Expires) {
		return nil, s.setHeader(errAccessTokenExpired, nil, w)
	}

	if len(scopes) == 0 {
		return at, nil
	}

	if len(scopes) > 0 && len(at.Scope) == 0 {
		return nil, s.setHeader(errInsufficientScope, scopes, w)
	}

	for i := 0; i < len(scopes); i++ {
		if checkScope(at.Scope, scopes[i]) {
			return at, nil
		}
	}

	return nil, s.setHeader(errInsufficientScope, scopes, w)
}

func (s *ResourceServer) setHeader(resp *Response, scopes []string, w http.ResponseWriter) *Response {
	strs := []string{}
	// authHeader := string2.Concat(`Bearer realm="`, s.config.WWWRealm, `"`)
	strs = append(strs, s.config.TokenBearerHeaderName, ` realm="`, s.config.WWWRealm, `"`)

	if len(scopes) > 0 {
		// authHeader = string2.Concat(authHeader, `, scopes="`, strings.Join(scopes, " "), `"`)
		strs = append(strs, `, scopes="`, strings.Join(scopes, " "), `"`)
	}

	if resp.IsError() {
		strs = append(strs, `, error="`, resp.data["error"].(string), `", error_description="`, resp.data["error_description"].(string), `"`)
		// authHeader = string2.Concat(authHeader, `, error="`, resp.data["error"].(string), `", error_description="`, resp.data["error_description"].(string), `"`)
	}

	w.Header().Set("WWW-Authenticate", string2.Concat(strs...))

	return resp
}

func (s *ResourceServer) getTokenFromHttp(r *http.Request) (string, *Response) {
	auth := r.Header.Get(`Authorization`)
	postAuth := r.PostFormValue(s.config.TokenParamName)
	getAuth := r.URL.Query().Get(s.config.TokenParamName)

	methodsUsed := 0
	if auth != "" {
		methodsUsed++
	}

	if getAuth != "" {
		methodsUsed++
	}

	if postAuth != "" {
		methodsUsed++
	}

	if methodsUsed > 1 {
		return "", errOnlyOneTokenMethod
	}

	if methodsUsed == 0 {
		return "", errNoTokenInRequest
	}

	if auth != "" {
		strs := strings.Fields(auth)
		if len(strs) < 2 || strs[0] != s.config.TokenBearerHeaderName {
			return "", errMalFormedHeader
		}
		return strs[1], nil
	}

	if postAuth != "" {
		return postAuth, nil
	}

	return getAuth, nil
}
