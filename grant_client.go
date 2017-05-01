package clover

type ClientCredentialsGrantType struct {
}

func NewClientCredentials() *ClientCredentialsGrantType {
	return &ClientCredentialsGrantType{}
}

func (g *ClientCredentialsGrantType) Validate(ctx Context) error {
	return nil
}

func (g *ClientCredentialsGrantType) Name() string {
	return "client_credentials"
}
