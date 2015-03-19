package clover

type authCodeGrant struct {
	store AuthServerStore
}

func newAuthCodeGrant(store AuthServerStore) GrantType {
	return &authCodeGrant{store}
}

func (g *authCodeGrant) Validate(tr *TokenRequest) (*GrantData, *response) {
	if tr.Code == "" {
		return nil, errCodeRequired
	}

	auth, err := g.store.GetAuthorizeCode(tr.Code)
	if err != nil {
		return nil, errAuthCodeNotExist
	}

	if tr.RedirectURI == "" || tr.RedirectURI != auth.RedirectURI {
		return nil, errRedirectMismatch
	}

	if isExpireUnix(auth.Expires) {
		return nil, errAuthCodeExpired
	}

	return &GrantData{
		ClientID: auth.ClientID,
		UserID:   auth.UserID,
		Scope:    auth.Scope,
	}, nil
}

func (g *authCodeGrant) GetGrantType() string {
	return AUTHORIZATION_CODE
}

func (g *authCodeGrant) CreateAccessToken(td *TokenData, respType ResponseType) *response {
	return respType.GetAccessToken(td, true)
}