package clover

type refreshGrant struct {
	store RefreshTokenStore
}

func NewRefreshToken(store RefreshTokenStore) GrantType {
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
		ClientID: rt.ClientID,
		UserID:   rt.UserID,
		Scope:    rt.Scope,
		Data:     rt.Data,
	}, nil
}

func (g *refreshGrant) Name() string {
	return REFRESH_TOKEN
}

func (g *refreshGrant) IncludeRefreshToken() bool {
	return true
}

func (g *refreshGrant) BeforeCreateAccessToken(tr *TokenRequest, td *TokenData) *Response {
	if err := g.store.RemoveRefreshToken(tr.RefreshToken); err != nil {
		return errInternal(err.Error())
	}

	return nil
}
