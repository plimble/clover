package clover

type AuthCodeReq struct {
}

type AuthCodeRes struct {
	Headers     map[string]string
	RedirectURI string
}

func (o *oauth2) AuthorizationCodeHandler(req *AuthCodeReq) (*AuthCodeRes, error) {
	return nil, ErrNotImplemented
}
