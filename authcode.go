package clover

type authCodeGrant struct {
}

func newAuthCodeGrant() GrantType {
	return &authCodeGrant{}
}

func (g *authCodeGrant) Validate(tr *tokenRequest, c *Clover) (*GrantData, *Response) {
	if tr.code == "" {
		return nil, errCodeRequired
	}

	auth, err := c.Config.Store.GetAuthorizeCode(tr.code)
	if err != nil {
		return nil, errAuthCodeNotExist
	}

	if tr.redirectURI == "" || tr.redirectURI != auth.RedirectURI {
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

func (g *authCodeGrant) CreateAccessToken(td *TokenData, c *Clover, respType ResponseType) *Response {
	return respType.GetAccessToken(td, c, true)
}
