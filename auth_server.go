package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type AuthServerConfig struct {
	StateParamRequired   bool
	DefaultScopes        []string
	AllowCredentialsBody bool
}

func DefaultAuthServerConfig() *AuthServerConfig {
	return &AuthServerConfig{
		StateParamRequired:   true,
		AllowCredentialsBody: false,
	}
}

type authorizeRequest struct {
	state        string
	redirectURI  string
	responseType string
	clientID     string
	scope        string
}

type TokenRequest struct {
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
}

type AuthServer struct {
	config        *AuthServerConfig
	store         ClientStore
	authRespTypes map[string]AuthorizeRespType
	tokenRespType AccessTokenRespType
	grants        map[string]GrantType
	authorizeCtrl *authorizeCtrl
	tokenCtrl     *tokenCtrl
}

func NewAuthServer(store ClientStore, config *AuthServerConfig) *AuthServer {
	a := &AuthServer{
		config:        config,
		store:         store,
		authRespTypes: make(map[string]AuthorizeRespType),
		// tokenRespType: ,
		grants: make(map[string]GrantType),
	}

	a.authorizeCtrl = newAuthorizeCtrl(a.store, a.config)
	a.tokenCtrl = newTokenCtrl(a.store, a.config)

	return a
}

func (a *AuthServer) AddGrantType(grantType GrantType) {
	a.grants[grantType.Name()] = grantType
}

func (a *AuthServer) AddRespType(authorizeRespType AuthorizeRespType) {
	a.authRespTypes[authorizeRespType.Name()] = authorizeRespType
}

func (a *AuthServer) SetAccessTokenRespType(accessTokenRespType AccessTokenRespType) {
	a.tokenRespType = accessTokenRespType
}

func (a *AuthServer) Authorize(w http.ResponseWriter, r *http.Request, isAuthorized bool, userID string) *Response {
	ar := parseAuthRequest(r)

	return a.authorizeCtrl.authorize(ar, a.authRespTypes, isAuthorized, userID)
}

func (a *AuthServer) ValidateAuthorize(w http.ResponseWriter, r *http.Request) (*AuthorizeData, *Response) {
	ar := parseAuthRequest(r)

	ad, resp := a.authorizeCtrl.validate(ar, a.authRespTypes)
	if resp != nil {
		return nil, resp
	}

	return ad, nil
}

func (a *AuthServer) Token(w http.ResponseWriter, r *http.Request) *Response {
	if strings.ToLower(r.Method) != "post" {
		return errMustBePostMetthod
	}

	tr, resp := parseTokenRequest(r, a.config)
	if resp != nil {
		return resp
	}

	return a.tokenCtrl.token(tr, a.tokenRespType, a.grants)
}

func parseTokenRequest(r *http.Request, config *AuthServerConfig) (*TokenRequest, *Response) {
	clientID, clientSecret, resp := getCredentialsFromHttp(r, config)
	if resp != nil {
		return nil, resp
	}

	tr := &TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        reqVal(r, "scope", false),
		GrantType:    reqVal(r, "grant_type", false),
		Username:     reqVal(r, "username", false),
		Password:     reqVal(r, "password", false),
		Code:         reqVal(r, "code", false),
		RedirectURI:  reqVal(r, "redirect_uri", false),
		RefreshToken: reqVal(r, "refresh_token", false),
		Assertion:    reqVal(r, "assertion", false),
	}

	return tr, nil
}

func parseAuthRequest(r *http.Request) *authorizeRequest {
	return &authorizeRequest{
		state:        reqVal(r, "state", true),
		redirectURI:  reqVal(r, "redirect_uri", true),
		responseType: reqVal(r, "response_type", true),
		clientID:     reqVal(r, "client_id", true),
		scope:        reqVal(r, "scope", true),
	}
}

func getCredentialsFromHttp(r *http.Request, config *AuthServerConfig) (string, string, *Response) {
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

func reqVal(r *http.Request, key string, allowQuery bool) string {
	val := r.PostFormValue(key)
	if allowQuery && val == "" {
		val = r.FormValue(key)
	}

	return val
}
