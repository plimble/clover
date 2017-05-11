package clover

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type HTTPContext struct {
	AuthorizationHeader string
	Method              string
	Form                url.Values
	request             *http.Request
	response            http.ResponseWriter
}

type AccessTokenContext struct {
	HTTPContext
	Client               Client
	RefreshToken         string
	Username             string
	Password             string
	Code                 string
	RedirectURI          string
	GrantType            string
	UserID               string
	Scopes               []string
	AccessTokenLifespan  int
	RefreshTokenLifespan int
	Extra                map[string]interface{}
}

type AuthorizeContext struct {
	HTTPContext
	Client       Client
	Challenge    string
	ClientID     string
	UserID       string
	Scopes       []string
	Extra        map[string]interface{}
	ResponseType string
	State        string
	RedirectURI  string
}

type IntrospectionContext struct {
	HTTPContext
	Client    Client
	Token     string
	TokenType string
}

type RevokeContext struct {
	HTTPContext
	Client    Client
	Token     string
	TokenType string
}

func parseHTTPRequest(w http.ResponseWriter, r *http.Request) (HTTPContext, error) {
	if err := r.ParseForm(); err != nil {
		return HTTPContext{}, errors.WithStack(errInvalidParseForm)
	}

	return HTTPContext{
		Method:              r.Method,
		AuthorizationHeader: r.Header.Get("Authorization"),
		Form:                r.Form,
		request:             r,
		response:            w,
	}, nil
}

func parseAccessTokenRequest(w http.ResponseWriter, r *http.Request) (*AccessTokenContext, error) {
	httpctx, err := parseHTTPRequest(w, r)
	if err != nil {
		return nil, err
	}

	return &AccessTokenContext{
		HTTPContext:  httpctx,
		GrantType:    httpctx.Form.Get("grant_type"),
		Scopes:       strings.Fields(httpctx.Form.Get("scope")),
		Code:         httpctx.Form.Get("code"),
		RedirectURI:  httpctx.Form.Get("redirect_uri"),
		Username:     httpctx.Form.Get("username"),
		Password:     httpctx.Form.Get("password"),
		RefreshToken: httpctx.Form.Get("refresh_token"),
	}, nil
}

func parseAuthorizeRequest(w http.ResponseWriter, r *http.Request) (*AuthorizeContext, error) {
	httpctx, err := parseHTTPRequest(w, r)
	if err != nil {
		return nil, err
	}

	ctx := &AuthorizeContext{
		HTTPContext: httpctx,
		Challenge:   httpctx.Form.Get("challenge"),
	}

	if ctx.Challenge == "" {
		ctx.ResponseType = httpctx.Form.Get("response_type")
		ctx.State = httpctx.Form.Get("state")
		ctx.ClientID = httpctx.Form.Get("client_id")
		ctx.RedirectURI = httpctx.Form.Get("redirect_uri")
		ctx.Scopes = strings.Fields(httpctx.Form.Get("scope"))
	}

	return ctx, nil
}

func parseIntospectionRequest(w http.ResponseWriter, r *http.Request) (*IntrospectionContext, error) {
	httpctx, err := parseHTTPRequest(w, r)
	if err != nil {
		return nil, err
	}

	return &IntrospectionContext{
		HTTPContext: httpctx,
		Token:       httpctx.Form.Get("token"),
		TokenType:   httpctx.Form.Get("token_type_hint"),
	}, nil
}

func parseRevokeRequest(w http.ResponseWriter, r *http.Request) (*RevokeContext, error) {
	httpctx, err := parseHTTPRequest(w, r)
	if err != nil {
		return nil, err
	}

	return &RevokeContext{
		HTTPContext: httpctx,
		Token:       httpctx.Form.Get("token"),
		TokenType:   httpctx.Form.Get("token_type_hint"),
	}, nil
}
