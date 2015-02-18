package clover

type passwordGrant struct {
}

func newPasswordGrant() GrantType {
	return &passwordGrant{}
}

func (g *passwordGrant) Validate(tr *tokenRequest, c *Clover) (*GrantData, *Response) {
	if tr.username == "" || tr.password == "" {
		return nil, errUsernamePasswordRequired
	}

	uid, err := c.Config.Store.GetUser(tr.username, tr.password)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	client, err := c.Config.Store.GetClient(tr.clientID)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	if client.ClientSecret != tr.clientSecret {
		return nil, errInvalidClientCredentail
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

func (g *passwordGrant) CreateAccessToken(td *TokenData, c *Clover, respType ResponseType) *Response {
	return respType.GetAccessToken(td, c, true)
}
