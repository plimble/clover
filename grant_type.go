package clover

import "github.com/pkg/errors"

type GrantData struct {
	UserID string
	Scopes []string
}

type GrantType interface {
	Validate(req *AccessTokenReq) error
	Name() string
	CreateAccessToken(req *AccessTokenReq) (*AccessTokenRes, error)
}

func DefaultGrantCheckScope(grantScopes []string, clientScopes []string) ([]string, error) {
	if len(grantScopes) > 0 {
		if len(clientScopes) > 0 {
			if ok := CheckScope(clientScopes, grantScopes); !ok {
				return nil, errors.WithStack(ErrInvalidScope)
			}
			return grantScopes, nil
		} else {
			return nil, errors.WithStack(ErrUnSupportedScope)
		}
	} else if len(clientScopes) > 0 {
		return clientScopes, nil
	}

	return nil, errors.WithStack(ErrUnSupportedScope)
}
