package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func addSecondUnix(sec int) int64 {
	if sec == 0 {
		return 0
	}

	return time.Now().UTC().Truncate(time.Nanosecond).Add(time.Second * time.Duration(sec)).Unix()
}

func isExpireUnix(expires int64) bool {
	return time.Now().UTC().Truncate(time.Nanosecond).Unix() > expires
}

func getCredentialsFromHttp(req *http.Request) (string, string, error) {
	headerAuth := req.Header.Get("Authorization")

	if headerAuth == "" {
		return "", "", errors.WithStack(ErrClientCredentialRequired)
	}

	s := strings.SplitN(headerAuth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", "", errors.WithStack(ErrInvalidAuthHeader)
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", ErrInternalServer.WithCause(errors.WithStack(err))
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", errors.WithStack(ErrInvalidAuthMessage)
	}

	return pair[0], pair[1], nil
}
