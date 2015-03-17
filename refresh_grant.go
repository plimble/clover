package clover

type refreshGrant struct {
}

func newRefreshGrant() GrantType {
	return &refreshGrant{}
}

func (g *refreshGrant) Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *response) {
	if tr.RefreshToken == "" {
		return nil, errRefreshTokenRequired
	}

	rt, err := a.Config.Store.GetRefreshToken(tr.RefreshToken)
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
	}, nil
}

func (g *refreshGrant) GetGrantType() string {
	return REFRESH_TOKEN
}

func (g *refreshGrant) CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *response {
	if err := a.Config.Store.RemoveRefreshToken(td.GrantData.RefreshToken); err != nil {
		return errInternal(err.Error())
	}

	return respType.GetAccessToken(td, a, true)
}
