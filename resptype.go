package clover

const (
	RESP_TYPE_CODE  = "code"
	RESP_TYPE_TOKEN = "token"
	RESP_TPYE_JWT   = "jwt"
)

// type ResponseType interface {
// 	GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *response
// 	GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *response
// 	GetResponseType() string
// }

type ResponseType interface {
	AuthResponseType
	AccessTokenResponseType
}

type AuthResponseType interface {
	GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *response
}

type AccessTokenResponseType interface {
	GetAccessToken(td *TokenData, includeRefresh bool) *response
}
