package clover

type passwordGrant struct {
	store AuthServerStore
}

func newPasswordGrant(store AuthServerStore) GrantType {
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

func (g *passwordGrant) CreateAccessToken(td *TokenData, respType ResponseType) *response {
	return respType.GetAccessToken(td, true)
}
