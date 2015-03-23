package clover

const (
	RESP_TYPE_CODE  = "code"
	RESP_TYPE_TOKEN = "token"
	RESP_TPYE_JWT   = "jwt"
)

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
