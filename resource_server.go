package clover

import (
	"fmt"
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
	token, resp := getTokenFromHttp(r)
	if resp != nil {
		return nil, s.responseError(resp, scopes, w)
	}

	at, err := s.config.AuthServerStore.GetAccessToken(token)
	if err != nil {
		return nil, s.responseError(errInvalidAccessToken, scopes, w)
	}

	if at.Expires > 0 && isExpireUnix(at.Expires) {
		return nil, s.responseError(errAccessTokenExpired, scopes, w)
	}

	if len(scopes) == 0 {
		return at, newRespData(nil)
	}

	if len(scopes) > 0 && len(at.Scope) == 0 {
		return nil, s.responseError(errInsufficientScope, scopes, w)
	}

	for _, scope := range scopes {
		if checkScope(at.Scope, scope) {
			return at, newRespData(nil)
		}
	}

	return nil, s.responseError(errInsufficientScope, scopes, w)
}

func (s *ResourceServer) responseError(resp *Response, scopes []string, w http.ResponseWriter) *Response {
	resp.setHeader(map[string]string{
		"WWW-Authenticate": fmt.Sprintf(`%s realm="%s", scope="%s", error="%s", error_description="%s"`,
			"Bearer", "Service", strings.Join(scopes, " "), resp.data["error"], resp.data["error_description"],
		),
	})
	return resp
}

func getTokenFromHttp(r *http.Request) (string, *Response) {
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
