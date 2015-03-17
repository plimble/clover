package clover

type ResponseType interface {
	GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *Response
	GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *Response
	GetResponseType() string
}
