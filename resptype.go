package clover

//go:generate mockery -name AuthorizeRespType -inpkg
//go:generate mockery -name AccessTokenRespType -inpkg

type AuthorizeRespType interface {
	Name() string
	Response(ad *AuthorizeData, userID string) *Response
	SupportGrant() string
	IsImplicit() bool
}

type AccessTokenRespType interface {
	Response(td *TokenData, includeRefresh bool) *Response
}
