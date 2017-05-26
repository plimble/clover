package token

import (
	"time"

	"github.com/plimble/clover/oauth2"
)

type ClientCredentialsGrantType struct {
	AccessTokenLifespan int
}

func (g *ClientCredentialsGrantType) GrantRequest(req *TokenHandlerRequest, client *oauth2.Client, storage oauth2.Storage) (*GrantData, error) {
	if client.Public {
		return nil, InvalidClient("public client is not allowed for client_credential grant type")
	}

	return &GrantData{
		Scopes:              client.Scopes,
		AccessTokenLifespan: g.AccessTokenLifespan,
		IncludeRefreshToken: false,
	}, nil
}

func (g *ClientCredentialsGrantType) Name() string {
	return "client_credentials"
}

func (g *ClientCredentialsGrantType) CreateToken(grantData *GrantData, client *oauth2.Client, storage oauth2.Storage, tokenGen oauth2.TokenGenerator) (string, string, error) {
	atoken, err := tokenGen.CreateAccessToken(&oauth2.CreateAccessTokenRequest{
		ClientID:  client.ID,
		Scopes:    grantData.Scopes,
		ExpiresIn: grantData.AccessTokenLifespan,
		Extras:    grantData.Extras,
	})
	if err != nil {
		return "", "", err
	}

	at := &oauth2.AccessToken{
		AccessToken: atoken,
		ClientID:    client.ID,
		Scopes:      grantData.Scopes,
		Expired:     time.Now().UTC().Add(time.Second * time.Duration(grantData.AccessTokenLifespan)).Unix(),
		ExpiresIn:   grantData.AccessTokenLifespan,
		Extras:      grantData.Extras,
	}

	if err = storage.SaveAccessToken(at); err != nil {
		return "", "", err
	}

	return atoken, "", nil
}
