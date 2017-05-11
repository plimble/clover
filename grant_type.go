package clover

// GrantType use this interface for implement extension grant type
type GrantType interface {
	Validate(ctx *AccessTokenContext, tokenManager TokenManager) error
	Name() string
	CreateAccessToken(ctx *AccessTokenContext, tokenManager TokenManager) (*AccessTokenRes, error)
}
