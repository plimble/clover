package clover

type authCodeGrant struct {
}

func newAuthCodeGrant() GrantType {
	return &authCodeGrant{}
}

func (g *authCodeGrant) Validate(tr *TokenRequest, a *AuthorizeServer) (*GrantData, *Response) {
	if tr.Code == "" {
		return nil, errCodeRequired
	}

	auth, err := a.Config.Store.GetAuthorizeCode(tr.Code)
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

func (g *authCodeGrant) CreateAccessToken(td *TokenData, a *AuthorizeServer, respType ResponseType) *Response {
	return respType.GetAccessToken(td, a, true)
}
