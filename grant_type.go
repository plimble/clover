package clover

const (
	AUTHORIZATION_CODE = "authorization_code"
	REFRESH_TOKEN      = "refresh_token"
	PASSWORD           = "password"
	CLIENT_CREDENTIALS = "client_credentials"
	IMPLICIT           = "implicit"
	// JWT_GRANT          = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

type Grants map[string]GrantType

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
	Validate(tr *TokenRequest) (*GrantData, *Response)
	GetGrantType() string
	CreateAccessToken(td *TokenData, respType TokenRespType) *Response
}
