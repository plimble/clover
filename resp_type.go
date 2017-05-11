package clover

import (
	"net/url"
)

type ResponseType interface {
	Name() string
	GenerateUrl(ctx *AuthorizeContext, tokenManager TokenManager) (*url.URL, error)
}
