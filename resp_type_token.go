package clover

import (
	"net/url"

	"github.com/pkg/errors"
)

type TokenResponseType struct {
	AccessTokenLifeSpan int
}

func (r *TokenResponseType) Name() string {
	return "token"
}

func (r *TokenResponseType) GenerateUrl(ctx *AuthorizeContext, tokenManager TokenManager) (*url.URL, error) {
	redirect, err := url.Parse(ctx.RedirectURI)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	acToken := &AccessTokenContext{
		Client:              ctx.Client,
		UserID:              ctx.UserID,
		Scopes:              ctx.Scopes,
		AccessTokenLifespan: r.AccessTokenLifeSpan,
	}

	accessToken, _, err := tokenManager.GenerateAccessToken(acToken, false)
	if err != nil {
		return nil, err
	}

	q := redirect.Query()
	q.Add("token", accessToken.AccessToken)
	q.Add("state", ctx.State)

	redirect.Fragment = q.Encode()

	return redirect, nil
}
