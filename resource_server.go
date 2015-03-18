package clover

import (
	"fmt"
	"net/http"
	"strings"
)

type ValidAccessToken func(at *AccessToken)

type ResourceServer struct {
	store Store
}

func NewResourceServer(store Store) *ResourceServer {
	return &ResourceServer{store}
}

func (s *ResourceServer) VerifyAccessToken(w http.ResponseWriter, r *http.Request, scopes []string, fn ValidAccessToken) {
	token, resp := getTokenFromHttp(r)
	if resp != nil {
		s.responseError(resp, scopes, w)
		return
	}

	at, err := s.store.GetAccessToken(token)
	if err != nil {
		s.responseError(errInvalidAccessToken, scopes, w)
		return
	}

	if at.Expires > 0 && isExpireUnix(at.Expires) {
		s.responseError(errAccessTokenExpired, scopes, w)
		return
	}

	if len(scopes) > 0 && len(at.Scope) == 0 {
		s.responseError(errInsufficientScope, scopes, w)
		return
	}

	for _, scope := range scopes {
		if checkScope(at.Scope, scope) {
			fn(at)
			return
		}
	}

	s.responseError(errInsufficientScope, scopes, w)
	return
}

func (s *ResourceServer) responseError(resp *response, scopes []string, w http.ResponseWriter) {
	resp.SetHeader(map[string]string{
		"WWW-Authenticate": fmt.Sprintf(`%s realm="%s", scope="%s", error="%s", error_description="%s"`,
			"Bearer", "Service", strings.Join(scopes, " "), resp.data["error"], resp.data["error_description"],
		),
	})
	resp.Write(w)
}

func getTokenFromHttp(r *http.Request) (string, *response) {
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
