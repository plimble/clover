package clover

import (
	"github.com/pkg/errors"
)

type AccessTokenRes struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       string `json:"user_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type accessTokenFlow struct {
	tokenManager   TokenManager
	grantTypes     map[string]GrantType
	scopeValidator ScopeValidator
	clientStorage  ClientStorage
}

func (f *accessTokenFlow) run(ctx *AccessTokenContext) (*AccessTokenRes, error) {
	if ctx.Method != "POST" {
		return nil, errors.WithStack(errMethodPostRequired)
	}

	if ctx.GrantType == "" {
		return nil, errors.WithStack(errGrantTypeNotFound)
	}

	grant, ok := f.grantTypes[ctx.GrantType]
	if !ok {
		return nil, errors.WithStack(ErrGrantTypeNotSupport(grant.Name()))
	}

	clientID, clientSecret, err := getCredentialsFromHttp(ctx.AuthorizationHeader)
	if err != nil {
		return nil, err
	}

	client, err := f.clientStorage.GetClientWithSecret(clientID, clientSecret)
	if err != nil {
		return nil, errInvalidClient.WithCause(err)
	}
	ctx.Client = *client

	if err = grant.Validate(ctx, f.tokenManager); err != nil {
		return nil, err
	}

	if f.scopeValidator != nil && grant.Name() != "client_credentials" && grant.Name() != "authorization_code" {
		ctx.Scopes, err = f.scopeValidator.Validate(ctx.Scopes, ctx.Client.Scopes)
		if err != nil {
			return nil, errInvalidScope.WithCause(err)
		}
	}

	res, err := grant.CreateAccessToken(ctx, f.tokenManager)
	if err != nil {
		return nil, errInternalServer.WithCause(err)
	}

	return res, nil
}
