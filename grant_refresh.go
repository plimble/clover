package clover

type refreshGrant struct {
	store RefreshTokenStore
}

func newRefreshGrant(store RefreshTokenStore) GrantType {
	return &refreshGrant{store}
}

func (g *refreshGrant) Validate(tr *TokenRequest) (*GrantData, *Response) {
	if tr.RefreshToken == "" {
		return nil, errRefreshTokenRequired
	}

	rt, err := g.store.GetRefreshToken(tr.RefreshToken)
	if err != nil {
		return nil, errInvalidRefreshToken
	}

	if isExpireUnix(rt.Expires) {
		return nil, errRefreshTokenExpired
	}

	return &GrantData{
		ClientID:     rt.ClientID,
		UserID:       rt.UserID,
		Scope:        rt.Scope,
		RefreshToken: rt.RefreshToken,
		Data:         rt.Data,
	}, nil
}

func (g *refreshGrant) GetGrantType() string {
	return REFRESH_TOKEN
}

func (g *refreshGrant) CreateAccessToken(td *TokenData, respType TokenRespType) *Response {
	if err := g.store.RemoveRefreshToken(td.GrantData.RefreshToken); err != nil {
		return errInternal(err.Error())
	}

	return respType.GetAccessToken(td, true)
}
