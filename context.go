package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type Context struct {
	Username     string
	Password     string
	Code         string
	RedirectURI  string
	RefreshToken string
	ClientID     string
	ClientSecret string
	GrantType    string
	Scope        string
	Assertion    string
	ResponseType string
	State        string
}

func ParseRequest(r *http.Request) (*Context, error) {
	var err error

	if err = r.ParseForm(); err != nil {
		return nil, ErrInvalidParseForm(err)
	}

	clientID, clientSecret, err := getCredentialsFromHttp(r)
	if err != nil {
		return nil, err
	}

	return &Context{
		Username:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Code:         r.FormValue("code"),
		RedirectURI:  r.FormValue("redirect_uri"),
		RefreshToken: r.FormValue("refresh_token"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    r.FormValue("grant_type"),
		Scope:        r.FormValue("scope"),
		Assertion:    r.FormValue("assertion"),
		ResponseType: r.FormValue("response_type"),
		State:        r.FormValue("state"),
	}, nil
}

func getCredentialsFromHttp(r *http.Request) (string, string, error) {
	headerAuth := r.Header.Get("Authorization")

	if headerAuth == "" {
		return "", "", ErrClientCredentialRequired()
	}

	s := strings.SplitN(headerAuth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", "", ErrInvalidAuthHeader()
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", ErrInternalServer(err)
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", "", ErrInvalidAuthMessage()
	}

	return pair[0], pair[1], nil
}
