package clover

type ResponseType interface {
	GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *response
	GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *response
	GetResponseType() string
}
