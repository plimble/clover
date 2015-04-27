package clover

type passwordGrant struct {
	store UserStore
}

func NewPassword(store UserStore) GrantType {
	return &passwordGrant{store}
}

func (g *passwordGrant) Validate(tr *TokenRequest) (*GrantData, *Response) {
	if tr.Username == "" || tr.Password == "" {
		return nil, errUsernamePasswordRequired
	}

	u, err := g.store.GetUser(tr.Username, tr.Password)

	if err != nil {
		return nil, errInvalidUsernamePAssword
	}

	return &GrantData{
		UserID: u.GetID(),
		Scope:  u.GetScope(),
		Data:   u.GetData(),
	}, nil
}

func (g *passwordGrant) Name() string {
	return PASSWORD
}

func (g *passwordGrant) IncludeRefreshToken() bool {
	return true
}

func (g *passwordGrant) BeforeCreateAccessToken(tr *TokenRequest, td *TokenData) *Response {
	return nil
}
