package clover

type refreshGrant struct {
}

func newRefreshGrant() GrantType {
	return &refreshGrant{}
}

func (g *refreshGrant) Validate(tr *tokenRequest, c *Clover) (*GrantData, *Response) {
	if tr.refreshToken == "" {
		return nil, errRefreshTokenRequired
	}

	rt, err := c.Config.Store.GetRefreshToken(tr.refreshToken)
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

func (g *refreshGrant) CreateAccessToken(td *TokenData, c *Clover, respType ResponseType) *Response {
	if err := c.Config.Store.RemoveRefreshToken(td.GrantData.RefreshToken); err != nil {
		return errInternal(err.Error())
	}

	return respType.GetAccessToken(td, c, true)
}
