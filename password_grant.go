package clover

type passwordGrant struct {
}

func newPasswordGrant() GrantType {
	return &passwordGrant{}
}

func (g *passwordGrant) Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *response) {
	if tr.Username == "" || tr.Password == "" {
		return nil, errUsernamePasswordRequired
	}

	uid, scopes, err := a.Config.Store.GetUser(tr.Username, tr.Password)
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

func (g *passwordGrant) CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *response {
	return respType.GetAccessToken(td, a, true)
}
