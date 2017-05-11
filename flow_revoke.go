package clover

import (
	"log"

	"github.com/pkg/errors"
)

type revokeFlow struct {
	clientStorage ClientStorage
	tokenManager  TokenManager
}

func (f *revokeFlow) Run(ctx *RevokeContext) error {
	token := ctx.Form.Get("token")
	tokenType := ctx.Form.Get("token_type_hint")

	if ctx.Method != "POST" {
		return errors.WithStack(errMethodPostRequired)
	}

	clientID, clientSecret, err := getCredentialsFromHttp(ctx.AuthorizationHeader)
	if err != nil {
		return err
	}

	_, err = f.clientStorage.GetClientWithSecret(clientID, clientSecret)
	if err != nil {
		return err
	}

	if tokenType != "access_token" && tokenType != "refresh_token" {
		return errors.WithStack(errRevokeTokenTypeInvalid)
	}

	if token == "" {
		return errors.WithStack(errRevokeTokenRequired)
	}

	switch tokenType {
	case "access_token":
		if err = f.tokenManager.DeleteAccessToken(token); err != nil {
			log.Printf("error: %+v", err)
			return nil
		}
	case "refresh_token":
		if err = f.tokenManager.DeleteRefreshToken(token); err != nil {
			log.Printf("error: %+v", err)
			return nil
		}
	}

	return nil
}
