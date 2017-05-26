package oauth2

import (
	"encoding/base64"
	"strings"
)

func GetCredentialsFromHttp(headerAuth string) (string, string, error) {
	if headerAuth == "" {
		return "", "", InvalidClient("client credentials are required")
	}

	s := strings.SplitN(headerAuth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", "", InvalidClient("invalid authorization header")
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", InvalidClient("cloud not decode Authorization").WithCause(err)
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", InvalidClient("invalid authorization message")
	}

	if pair[0] == "" || pair[1] == "" {
		return "", "", InvalidClient("client credentials are required")
	}

	return pair[0], pair[1], nil
}
