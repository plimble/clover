package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"
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

func getCredentialsFromHttp(r *http.Request, config *AuthorizeConfig) (string, string, *Response) {
	headerAuth := r.Header.Get("Authorization")

	switch {
	case headerAuth != "":
		s := strings.SplitN(headerAuth, " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			return "", "", errInvalidAuthHeader
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return "", "", errInternal(err.Error())
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			return "", "", errInvalidAuthMSG
		}

		return pair[0], pair[1], nil
	case config.AllowCredentialsBody:
		if r.PostForm.Get(`client_id`) == "" || r.PostForm.Get(`client_secret`) == "" {
			return "", "", errCredentailsNotInBody
		}

		return r.PostForm.Get("client_id"), r.PostForm.Get("client_secret"), nil
	}

	return "", "", errCredentailsRequired
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
