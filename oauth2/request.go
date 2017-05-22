package oauth2

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func GetAccessTokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		err := req.ParseForm()
		if err != nil {
			return ""
		}
		return req.Form.Get("access_token")
	}

	return split[1]
}

func GetCredentialsFromHttp(headerAuth string) (string, string, error) {
	if headerAuth == "" {
		return "", "", Error(400, "invalid_client", "client credentials are required")
	}

	s := strings.SplitN(headerAuth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", "", Error(400, "invalid_client", "invalid authorization header")
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", Error(400, "invalid_client", "cannot decode Authorization").WithCause(err)
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", Error(400, "invalid_client", "invalid authorization message")
	}

	if pair[0] == "" || pair[1] == "" {
		return "", "", Error(400, "invalid_client", "client credentials are required")
	}

	return pair[0], pair[1], nil
}
