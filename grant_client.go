package clover

type clientGrant struct {
	store ClientStore
}

func NewClientCredential(store ClientStore) GrantType {
	return &clientGrant{store}
}

func (g *clientGrant) Validate(tr *TokenRequest) (*GrantData, *Response) {
	client, err := g.store.GetClient(tr.ClientID)
	if err != nil {
		return nil, errInvalidClientCredentail
	}

	if client.GetClientSecret() != tr.ClientSecret {
		return nil, errInvalidClientCredentail
	}

	return &GrantData{
		ClientID: client.GetClientID(),
		UserID:   client.GetUserID(),
		Scope:    client.GetScope(),
		Data:     client.GetData(),
	}, nil
}

func (g *clientGrant) Name() string {
	return CLIENT_CREDENTIALS
}

func (g *clientGrant) IncludeRefreshToken() bool {
	return false
}

func (g *clientGrant) BeforeCreateAccessToken(tr *TokenRequest, td *TokenData) *Response {
	return nil
}
