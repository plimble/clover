package clover

type TokenType interface{}

type bearer struct {
}

func newBearerTokenType() *bearer {
	return &bearer{}
}
