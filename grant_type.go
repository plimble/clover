package clover

//go:generate mockery -name GrantType -inpkg

const (
	AUTHORIZATION_CODE = "authorization_code"
	REFRESH_TOKEN      = "refresh_token"
	PASSWORD           = "password"
	CLIENT_CREDENTIALS = "client_credentials"
	IMPLICIT           = "implicit"
	// JWT_GRANT          = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

type GrantData struct {
	ClientID string
	UserID   string
	Scope    []string
	Data     map[string]interface{}
}

type GrantType interface {
	Validate(tr *TokenRequest) (*GrantData, *Response)
	Name() string
	IncludeRefreshToken() bool
	BeforeCreateAccessToken(tr *TokenRequest, td *TokenData) *Response
}

func checkGrantType(available []string, request string) bool {
	for i := 0; i < len(available); i++ {
		if available[i] == request {
			return true
		}
	}

	return false
}
