package clover

import (
	"net/url"

	"github.com/pkg/errors"
)

type CodeResponseType struct {
	AuthorizeCodeLifespan int
}

func (r *CodeResponseType) Name() string {
	return "code"
}

func (r *CodeResponseType) GenerateUrl(ctx *AuthorizeContext, tokenManager TokenManager) (*url.URL, error) {
	redirect, err := url.Parse(ctx.RedirectURI)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	authCode, err := tokenManager.GenerateAuthorizeCode(ctx, r.AuthorizeCodeLifespan)
	if err != nil {
		return nil, err
	}

	q := redirect.Query()
	q.Add("code", authCode.Code)
	q.Add("state", ctx.State)

	redirect.RawQuery = q.Encode()

	return redirect, nil
}
