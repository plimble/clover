package clover

import (
	"net/http"

	"github.com/pkg/errors"
)

type RevokeReq struct {
	ClientID     string
	ClientSecret string
	Token        string
	TokenType    string
}

func NewRevokeReq(req *http.Request) (*RevokeReq, error) {
	var err error

	if err = req.ParseForm(); err != nil {
		return nil, errors.WithStack(ErrInvalidParseForm)
	}

	clientID, clientSecret, err := getCredentialsFromHttp(req)
	if err != nil {
		return nil, err
	}

	return &RevokeReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Token:        req.FormValue("token"),
		TokenType:    req.FormValue("token_type"),
	}, nil
}

func (o *oauth2) RevokeHandler(req *RevokeReq) error {
	if req.ClientID == "" || req.ClientSecret == "" {
		return errors.WithStack(ErrClientCredentialRequired)
	}

	var err error
	if _, err = o.clientManager.GetWithSecret(req.ClientID, req.ClientSecret); err != nil {
		return err
	}

	if req.TokenType != "access_token" && req.TokenType != "refresh_token" {
		return errors.WithStack(ErrRevokeTokenTypeInvalid)
	}

	if req.Token == "" {
		return errors.WithStack(ErrRevokeTokenRequired)
	}

	switch req.TokenType {
	case "access_token":
		if err = o.tokenManager.DeleteAccessToken(req.Token); err != nil {
			return err
		}
	case "refresh_token":
		if err = o.tokenManager.DeleteRefreshToken(req.Token); err != nil {
			return err
		}
	}

	return nil
}
