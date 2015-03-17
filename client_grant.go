package clover

type clientGrant struct {
}

func newClientGrant() GrantType {
	return &clientGrant{}
}

func (g *clientGrant) Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *Response) {
	client, err := a.Config.Store.GetClient(tr.ClientID)
	if err != nil {
		return nil, errInternal(err.Error())
	}

	if client.GetClientSecret() != tr.ClientSecret {
		return nil, errInvalidClientCredentail
	}

	return &GrantData{
		ClientID:  client.GetClientID(),
		UserID:    client.GetUserID(),
		Scope:     client.GetScope(),
		GrantType: client.GetGrantType(),
	}, nil
}

func (g *clientGrant) GetGrantType() string {
	return CLIENT_CREDENTIALS
}

func (g *clientGrant) CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *Response {
	return respType.GetAccessToken(td, a, false)
}
