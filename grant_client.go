package clover

type clientGrant struct {
	store AuthServerStore
}

func newClientGrant(store AuthServerStore) GrantType {
	return &clientGrant{store}
}

func (g *clientGrant) Validate(tr *TokenRequest) (*GrantData, *Response) {
	client, err := g.store.GetClient(tr.ClientID)
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

func (g *clientGrant) CreateAccessToken(td *TokenData, respType TokenRespType) *Response {
	return respType.GetAccessToken(td, false)
}
