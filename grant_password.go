package clover

type passwordGrant struct {
	store UserStore
}

func newPasswordGrant(store UserStore) GrantType {
	return &passwordGrant{store}
}

func (g *passwordGrant) Validate(tr *TokenRequest) (*GrantData, *response) {
	if tr.Username == "" || tr.Password == "" {
		return nil, errUsernamePasswordRequired
	}

	uid, scopes, err := g.store.GetUser(tr.Username, tr.Password)
	if err != nil {
		return nil, errInvalidUsernamePAssword
	}

	return &GrantData{
		ClientID: "",
		UserID:   uid,
		Scope:    scopes,
	}, nil
}

func (g *passwordGrant) GetGrantType() string {
	return PASSWORD
}

func (g *passwordGrant) CreateAccessToken(td *TokenData, respType AccessTokenResponseType) *response {
	return respType.GetAccessToken(td, true)
}
