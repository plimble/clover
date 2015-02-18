package clover

type ResponseType interface {
	GetAccessToken(tr *TokenData, c *Clover, includeRefresh bool) *Response
	GetAuthorizeResponse(ad *AuthorizeData, c *Clover) *Response
	GetResponseType() string
}
