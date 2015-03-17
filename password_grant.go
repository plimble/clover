package clover

type passwordGrant struct {
}

func newPasswordGrant() GrantType {
	return &passwordGrant{}
}

func (g *passwordGrant) Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *Response) {
	if tr.Username == "" || tr.Password == "" {
		return nil, errUsernamePasswordRequired
	}

	uid, err := a.Config.Store.GetUser(tr.Username, tr.Password)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	return &GrantData{
		ClientID: "",
		UserID:   uid,
		Scope:    []string{},
	}, nil
}

func (g *passwordGrant) GetGrantType() string {
	return PASSWORD
}

func (g *passwordGrant) CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *Response {
	return respType.GetAccessToken(td, a, true)
}
