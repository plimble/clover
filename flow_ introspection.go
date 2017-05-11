package clover

import (
	"strings"

	"github.com/pkg/errors"
)

type IntrospectionRes struct {
	Active       bool   `json:"active"`
	ClientID     string `json:"client_id"`
	Username     string `json:"username,omitempty"`
	Scopes       string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	// Exp          string `json:"exp,omitempty"`
	// Iat          string `json:"iat,omitempty"`
	// Sub          string `json:"sub,omitempty"`
	// Aud          string `json:"aud,omitempty"`
	// Iss          string `json:"iss,omitempty"`
	// Jti          string `json:"jti,omitempty"`
}

type introspectionFlow struct {
	clientStorage ClientStorage
	tokenManager  TokenManager
}

func (f *introspectionFlow) Run(ctx *IntrospectionContext) (*IntrospectionRes, error) {
	if ctx.Method != "POST" {
		return nil, errors.WithStack(errMethodPostRequired)
	}

	clientID, clientSecret, err := getCredentialsFromHttp(ctx.AuthorizationHeader)
	if err != nil {
		return nil, err
	}

	_, err = f.clientStorage.GetClientWithSecret(clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	if ctx.Token == "" {
		return nil, errors.WithStack(errIntroTokenRequired)
	}

	if ctx.TokenType == "" {
		return nil, errors.WithStack(errIntroTokenTypeRequired)
	}

	switch ctx.TokenType {
	case "access_token":
		accessToken, err := f.tokenManager.GetAccessToken(ctx.Token)
		if err != nil {
			return nil, err
		}

		return &IntrospectionRes{
			Active:    isExpireUnix(accessToken.Expired),
			ClientID:  accessToken.ClientID,
			Username:  accessToken.UserID,
			Scopes:    strings.Join(accessToken.Scopes, " "),
			TokenType: "access_token",
		}, nil
	case "refresh_token":
		refreshToken, err := f.tokenManager.GetRefreshToken(ctx.Token)
		if err != nil {
			return nil, err
		}

		return &IntrospectionRes{
			Active:    isExpireUnix(refreshToken.Expired),
			ClientID:  refreshToken.ClientID,
			TokenType: "refresh_token",
		}, nil
	}

	return nil, errors.WithStack(errIntroTokenTypeInvalid)
}
