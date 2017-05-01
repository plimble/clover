package clover

type GrantData struct {
	ClientID string
	UserID   string
	Scopes   []string
}

type GrantType interface {
	Validate(ctx Context) (*GrantData, error)
	Name() string
}
