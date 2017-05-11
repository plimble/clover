package clover

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func addSecondUnix(sec int) int64 {
	if sec == 0 {
		return 0
	}

	return time.Now().UTC().Truncate(time.Nanosecond).Add(time.Minute * time.Duration(sec)).Unix()
}

func isExpireUnix(expires int64) bool {
	return time.Now().UTC().Truncate(time.Nanosecond).Unix() > expires
}

func getCredentialsFromHttp(headerAuth string) (string, string, error) {
	if headerAuth == "" {
		return "", "", errors.WithStack(errClientCredentialRequired)
	}

	s := strings.SplitN(headerAuth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", "", errors.WithStack(errInvalidAuthHeader)
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", errInternalServer.WithCause(errors.WithStack(err))
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", errors.WithStack(errInvalidAuthMessage)
	}

	if pair[0] == "" || pair[1] == "" {
		return "", "", errors.WithStack(errClientCredentialRequired)
	}

	return pair[0], pair[1], nil
}

func getAccessTokenFromRequest(req *http.Request) string {
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

// RandomBytes returns n random bytes by reading from crypto/rand.Reader
func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}
