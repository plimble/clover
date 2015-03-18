package clover

const (
	AUTHORIZATION_CODE = "authorization_code"
	REFRESH_TOKEN      = "refresh_token"
	PASSWORD           = "password"
	CLIENT_CREDENTIALS = "client_credentials"
	IMPLICIT           = "implicit"
	JWT_GRANT          = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

type GrantData struct {
	ClientID     string
	UserID       string
	Scope        []string
	GrantType    []string
	RefreshToken string
}

func checkGrantType(grants []string, grant string) bool {
	for _, v := range grants {
		if grant == v {
			return true
		}
	}

	return false
}

type GrantType interface {
	Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *response)
	GetGrantType() string
	CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *response
}

func (a *AuthorizeServer) RegisterGrant(key string, grant GrantType) {
	a.Grant[key] = grant
}

func (a *AuthorizeServer) RegisterAuthCodeGrant() {
	a.Grant[AUTHORIZATION_CODE] = newAuthCodeGrant()
	a.RespType[RESP_TYPE_CODE] = newCodeResponseType()
}

func (a *AuthorizeServer) RegisterClientGrant() {
	a.Grant[REFRESH_TOKEN] = newRefreshGrant()
}

func (a *AuthorizeServer) RegisterPasswordGrant() {
	a.Grant[PASSWORD] = newPasswordGrant()
}

func (a *AuthorizeServer) RegisterRefreshGrant() {
	a.Grant[CLIENT_CREDENTIALS] = newClientGrant()
}

func (a *AuthorizeServer) RegisterImplicitGrant() {
	a.Grant[IMPLICIT] = newAuthCodeGrant()
}