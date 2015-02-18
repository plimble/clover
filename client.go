package clover

type clientGrant struct {
}

func newClientGrant() GrantType {
	return &clientGrant{}
}

func (g *clientGrant) Validate(tr *tokenRequest, c *Clover) (*GrantData, *Response) {
	client, err := c.Config.Store.GetClient(tr.clientID)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	if client.ClientSecret != tr.clientSecret {
		return nil, errInvalidClientCredentail
	}

	return &GrantData{
		ClientID:  client.ClientID,
		UserID:    client.UserID,
		Scope:     client.Scope,
		GrantType: client.GrantType,
	}, nil
}

func (g *clientGrant) GetGrantType() string {
	return CLIENT_CREDENTIALS
}

func (g *clientGrant) CreateAccessToken(td *TokenData, c *Clover, respType ResponseType) *Response {
	return respType.GetAccessToken(td, c, false)
}
