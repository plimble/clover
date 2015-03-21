package clover

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type Config struct {
	AccessLifeTime       int
	AuthCodeLifetime     int
	RefreshTokenLifetime int

	DefaultScopes        []string
	AllowCredentialsBody bool
	AllowImplicit        bool
	StateParamRequired   bool
}

type AuthorizeServer struct {
	stores        *stores
	Config        *Config
	tokenRespType ResponseType
	grant         map[string]GrantType
	authCtrl      *authController
	tokenCtrl     *tokenController
}

func NewAuthServer(store AuthServerStore, config *Config) *AuthorizeServer {
	a := &AuthorizeServer{
		stores: &stores{authServer: store},
		Config: config,
		grant:  make(map[string]GrantType),
	}

	return a
}

func DefaultConfig() *Config {
	return &Config{
		AccessLifeTime:       3600,
		AuthCodeLifetime:     60,
		RefreshTokenLifetime: 3600,
		AllowCredentialsBody: false,
		AllowImplicit:        false,
		StateParamRequired:   false,
	}
}

func (a *AuthorizeServer) UseJWTAccessTokens(store PublicKeyStore) {
	a.stores.publicKey = store
}

func (a *AuthorizeServer) RegisterGrant(name string, grant GrantType) {
	a.grant[name] = grant
}

func (a *AuthorizeServer) RegisterAuthCodeGrant(store AuthCodeStore) {
	a.stores.code = store
	a.grant[AUTHORIZATION_CODE] = newAuthCodeGrant(store)
}

func (a *AuthorizeServer) RegisterClientGrant() {
	a.grant[CLIENT_CREDENTIALS] = newClientGrant(a.stores.authServer)
}

func (a *AuthorizeServer) RegisterPasswordGrant(store UserStore) {
	a.stores.user = store
	a.grant[PASSWORD] = newPasswordGrant(store)
}

func (a *AuthorizeServer) RegisterRefreshGrant(store RefreshTokenStore) {
	a.stores.refresh = store
	a.grant[REFRESH_TOKEN] = newRefreshGrant(store)
}

func (a *AuthorizeServer) SetDefaultScopes(ids ...string) {
	a.Config.DefaultScopes = ids
}

func (a *AuthorizeServer) Authorize(w http.ResponseWriter, r *http.Request, isAuthorized bool) {
	ctrl := a.getAuthController()

	ar := a.parseAuthorizeRequest(r)
	resp := ctrl.handleAuthorize(ar, isAuthorized)

	resp.Write(w)
}

func (a *AuthorizeServer) ValidateAuthorize(w http.ResponseWriter, r *http.Request, fn ValidAuthorize) {
	ctrl := a.getAuthController()
	ar := a.parseAuthorizeRequest(r)

	client, scopes, _, resp := ctrl.validateAuthorizeRequest(ar)
	if resp != nil {
		resp.Write(w)
		return
	}

	fn(client, scopes)
}

func (a *AuthorizeServer) Token(w http.ResponseWriter, r *http.Request) {
	ctrl := a.getTokenController()
	if strings.ToLower(r.Method) != "post" {
		errMustBePostMetthod.Write(w)
		return
	}

	tr, resp := a.parseTokenRequest(r)
	if resp != nil {
		resp.Write(w)
		return
	}

	resp = ctrl.handleAccessToken(tr)
	resp.Write(w)
}

func (a *AuthorizeServer) parseTokenRequest(r *http.Request) (*TokenRequest, *response) {
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

func (a *AuthorizeServer) parseAuthorizeRequest(r *http.Request) *authorizeRequest {
	ar := &authorizeRequest{
		state:        r.FormValue("state"),
		redirectURI:  r.FormValue("redirect_uri"),
		responseType: strings.ToLower(r.FormValue("response_type")),
		clientID:     r.FormValue("client_id"),
		scope:        r.FormValue("scope"),
	}

	return ar
}

func (a *AuthorizeServer) getCredentialsFromHttp(r *http.Request, config *Config) (string, string, *response) {
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

func (a *AuthorizeServer) getAuthController() *authController {
	if a.authCtrl == nil {
		a.authCtrl = newAuthController(
			a.Config,
			a.stores.authServer,
			newAuthCodeResponseType(a.stores.code, a.Config),
			a.getRespType(),
		)
	}

	return a.authCtrl
}

func (a *AuthorizeServer) getTokenController() *tokenController {
	if a.tokenCtrl == nil {
		a.tokenCtrl = newTokenController(
			a.Config,
			a.stores.authServer,
			a.grant,
			a.getRespType(),
		)
	}

	return a.tokenCtrl
}

func (a *AuthorizeServer) getRespType() ResponseType {
	if a.tokenRespType == nil {
		if a.stores.publicKey != nil {
			a.tokenRespType = newJWTResponseType(a.stores.publicKey, a.stores.authServer, a.stores.refresh, a.Config)

		} else {
			a.tokenRespType = newTokenResponseType(a.stores.authServer, a.stores.refresh, a.Config)
		}
	}

	return a.tokenRespType
}
