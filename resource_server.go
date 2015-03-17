package clover

import (
	"fmt"
	"net/http"
	"strings"
)

type ResourceServer struct {
	store Store
}

func NewResourceServer(store Store) *ResourceServer {
	return &ResourceServer{store}
}

func (s *ResourceServer) VerifyResourceRequest(w http.ResponseWriter, r *http.Request, scope ...string) *AccessToken {
	token, resp := getTokenFromHttp(r)
	if resp != nil {
		resp.Write(w)
		return nil
	}

	at, err := s.store.GetAccessToken(token)
	if err != nil {
		errInvalidAccessToken.Write(w)
		return nil
	}

	if at.Expires > 0 && isExpireUnix(at.Expires) {
		errAccessTokenExpired.Write(w)
		return nil
	}

	if (len(scope) > 0 && len(at.Scope) == 0) || len(at.Scope) == 0 || checkScope(scope, at.Scope) {
		resp := errInsufficientScope
		resp.SetHeader(map[string]string{
			"WWW-Authenticate": fmt.Sprintf(`%s realm="%s", scope="%s", error="%s", error_description="%s`,
				"Bearer", "Service", strings.Join(scope, " "), resp.data["error"], resp.data["error_description"],
			),
		})
		resp.Write(w)
		return nil
	}

	return at
}
