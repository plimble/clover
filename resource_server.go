package clover

import (
	"github.com/plimble/utils/strings2"
	"net/http"
	"strings"
)

type ResourceServer struct {
	config *ResourceConfig
}

func NewResourceServer(config *ResourceConfig) *ResourceServer {
	return &ResourceServer{config}
}

func (s *ResourceServer) VerifyAccessToken(w http.ResponseWriter, r *http.Request, scopes ...string) (*AccessToken, *Response) {
	token, resp := s.getTokenFromHttp(r)
	if resp != nil {
		return nil, s.setHeader(resp, nil, w)
	}

	at, err := s.config.AuthServerStore.GetAccessToken(token)
	if err != nil {
		return nil, s.setHeader(errInvalidAccessToken, nil, w)
	}

	if at.Expires != 0 && isExpireUnix(at.Expires) {
		return nil, s.setHeader(errAccessTokenExpired, nil, w)
	}

	if len(scopes) == 0 {
		return at, s.setHeader(newRespData(nil), nil, w)
	}

	if len(scopes) > 0 && len(at.Scope) == 0 {
		return nil, s.setHeader(errInsufficientScope, scopes, w)
	}

	if !checkScope(at.Scope, scopes...) {
		return nil, s.setHeader(errInsufficientScope, scopes, w)
	}

	return at, s.setHeader(newRespData(nil), scopes, w)
}

func (s *ResourceServer) setHeader(resp *Response, scopes []string, w http.ResponseWriter) *Response {
	strs := []string{}
	// authHeader := string2.Concat(`Bearer realm="`, s.config.WWWRealm, `"`)
	strs = append(strs, `Bearer realm="`, s.config.WWWRealm, `"`)

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
	postAuth := r.PostFormValue("access_token")
	getAuth := r.URL.Query().Get("access_token")

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
		if len(strs) < 2 || strs[0] != "Bearer" {
			return "", errMalFormedHeader
		}
		return strs[1], nil
	}

	if postAuth != "" {
		return postAuth, nil
	}

	return getAuth, nil
}
