package clover

import (
	"net/http"
	"strings"

	"net/url"

	"github.com/pkg/errors"
)

type AccessTokenReq struct {
	Form         url.Values
	ClientID     string
	ClientSecret string
	Scopes       []string
	GrantType    string
	Session      string
	Client       *Client
	UserID       string
	data         map[string]interface{}
}

func NewAccessTokenReq(req *http.Request) (*AccessTokenReq, error) {
	var err error

	if err = req.ParseForm(); err != nil {
		return nil, errors.WithStack(ErrInvalidParseForm)
	}

	clientID, clientSecret, err := getCredentialsFromHttp(req)
	if err != nil {
		return nil, err
	}

	return &AccessTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    req.FormValue("grant_type"),
		Scopes:       strings.Fields(req.FormValue("scope")),
		Form:         req.Form,
		data:         make(map[string]interface{}),
	}, nil
}

func (req *AccessTokenReq) Get(key string) interface{} {
	return req.data[key]
}

func (req *AccessTokenReq) Set(key string, val interface{}) {
	req.data[key] = val
}

type AccessTokenRes struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       string `json:"user_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (o *oauth2) AccessTokenHandler(req *AccessTokenReq) (*AccessTokenRes, error) {
	if req.GrantType == "" {
		return nil, errors.WithStack(ErrGrantTypeNotFound)
	}

	if req.ClientID == "" || req.ClientSecret == "" {
		return nil, errors.WithStack(ErrClientCredentialRequired)
	}

	grant, ok := o.grantTypes[req.GrantType]
	if !ok {
		return nil, errors.WithStack(ErrGrantTypeNotSupport(grant.Name()))
	}

	var err error
	req.Client, err = o.clientManager.GetWithSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	if err = grant.Validate(req); err != nil {
		return nil, err
	}

	return grant.CreateAccessToken(req)
}
