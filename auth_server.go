package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type AuthServer struct {
	Config    *AuthConfig
	authCtrl  *authController
	tokenCtrl *tokenController
}

func NewAuthServer(config *AuthConfig) *AuthServer {
	a := &AuthServer{
		Config: config,
	}

	tokenResp := a.getTokenRespType()
	authResp := newAuthCodeResponseType(a.Config)

	a.tokenCtrl = newTokenController(a.Config, tokenResp)
	a.authCtrl = newAuthController(a.Config, authResp, tokenResp)

	return a
}

func (a *AuthServer) Authorize(w http.ResponseWriter, r *http.Request, isAuthorized bool) *Response {
	ar := a.parseAuthRequest(r)
	resp := a.authCtrl.handleAuthorize(ar, isAuthorized)

	return resp
}

func (a *AuthServer) ValidateAuthorize(w http.ResponseWriter, r *http.Request) (Client, []string, *Response) {
	ar := a.parseAuthRequest(r)

	client, scopes, _, resp := a.authCtrl.validateAuthorizeRequest(ar)
	if resp != nil {
		return nil, nil, resp
	}

	return client, scopes, newRespData(nil)
}

func (a *AuthServer) Token(w http.ResponseWriter, r *http.Request) *Response {
	if strings.ToLower(r.Method) != "post" {
		return errMustBePostMetthod
	}

	tr, resp := a.parseTokenRequest(r)
	if resp != nil {
		return resp
	}

	return a.tokenCtrl.handleAccessToken(tr)
}

func (a *AuthServer) parseTokenRequest(r *http.Request) (*TokenRequest, *Response) {
	clientID, clientSecret, resp := a.getCredentialsFromHttp(r, a.Config)
	if resp != nil {
		return nil, resp
	}

	tr := &TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        r.FormValue("scope"),
		GrantType:    r.FormValue("grant_type"),
		Username:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Code:         r.FormValue("code"),
		RedirectURI:  r.FormValue("redirect_uri"),
		RefreshToken: r.FormValue("refresh_token"),
		Assertion:    r.FormValue("assertion"),
	}

	return tr, nil
}

func (a *AuthServer) parseAuthRequest(r *http.Request) *authorizeRequest {
	ar := &authorizeRequest{
		state:        r.FormValue("state"),
		redirectURI:  r.FormValue("redirect_uri"),
		responseType: strings.ToLower(r.FormValue("response_type")),
		clientID:     r.FormValue("client_id"),
		scope:        r.FormValue("scope"),
	}

	return ar
}

func (a *AuthServer) getCredentialsFromHttp(r *http.Request, config *AuthConfig) (string, string, *Response) {
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
		if r.PostForm == nil {
			r.ParseForm()
		}
		if r.PostForm.Get(`client_id`) == "" || r.PostForm.Get(`client_secret`) == "" {
			return "", "", errCredentailsNotInBody
		}

		return r.PostForm.Get("client_id"), r.PostForm.Get("client_secret"), nil
	}

	return "", "", errCredentailsRequired
}

func (a *AuthServer) getTokenRespType() ResponseType {
	respType := newTokenRespType(a.Config)

	if a.Config.PublicKeyStore != nil {
		return newJWTResponseType(respType, a.Config)
	}

	return respType
}
