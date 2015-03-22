package clover

const (
	RESP_TYPE_CODE  = "code"
	RESP_TYPE_TOKEN = "token"
	RESP_TPYE_JWT   = "jwt"
)

// type ResponseType interface {
// 	GetAccessToken(td *TokenData, a *AuthorizeServer, includeRefresh bool) *Response
// 	GetAuthorizeResponse(client Client, scopes []string, ar *authorizeRequest, a *AuthorizeServer) *Response
// 	GetResponseType() string
// }

type ResponseType interface {
	AuthRespType
	TokenRespType
}

type AuthRespType interface {
	GetAuthResponse(ar *authorizeRequest, client Client, scopes []string) *Response
}

type TokenRespType interface {
	GetAccessToken(td *TokenData, includeRefresh bool) *Response
}
