package clover

type VerifyAccessTokenReq struct {
}

type VerifyAccessTokenRes struct{}

func (o *oauth2) VerifyAccessTokenHandler(req *VerifyAccessTokenReq) (*VerifyAccessTokenRes, error) {
	return nil, ErrNotImplemented
}
