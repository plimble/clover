package token

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type ProcessFunc func(req *TokenHandlerRequest) (TokenHandlerResponse, error)

type TokenHandlerRequest struct {
	ClientID     string
	ClientSecret string
	GrantType    string
	Form         url.Values
}

type TokenHandlerResponse map[string]interface{}

type TokenHandler struct {
	GrantTypes     map[string]GrantType
	Storage        oauth2.Storage
	TokenGenerator oauth2.TokenGenerator
	*zap.Logger
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		oauth2.WriteJsonError(w, ErrNotPostMethod)
		return
	}

	auth := r.Header.Get("Authorization")
	clientID, clientSecret, err := oauth2.GetCredentialsFromHttp(auth)
	if err != nil {
		oauth2.WriteJsonError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		oauth2.WriteJsonError(w, ErrParseForm)
		return
	}

	req := &TokenHandlerRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    r.Form.Get("grant_type"),
		Form:         r.Form,
	}

	res, err := h.RequestToken(req)
	h.Logger.Info("RequestToken",
		zap.Any("TokenHandlerRequest", req),
		zap.Any("TokenHandlerResponse", res),
		zap.Error(err),
	)
	if err != nil {
		oauth2.WriteJsonError(w, err)
		return
	}

	oauth2.WriteJson(w, 200, res)
}

func (h *TokenHandler) RequestToken(req *TokenHandlerRequest) (TokenHandlerResponse, error) {
	grant, err := h.getGrant(req)
	if err != nil {
		return nil, err
	}

	client, err := h.getClient(req)
	if err != nil {
		return nil, err
	}

	grantData, err := grant.GrantRequest(req, client, h.Storage)
	if err != nil {
		return nil, err
	}

	if err = h.checkScope(req.GrantType, grantData.Scopes, client.Scopes); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := grant.CreateToken(grantData, client, h.Storage, h.TokenGenerator)

	return h.createResponse(grantData, accessToken, refreshToken), err
}

func (h *TokenHandler) getGrant(req *TokenHandlerRequest) (GrantType, error) {
	if req.GrantType == "" {
		return nil, ErrGrantTypeRequired
	}

	grant, ok := h.GrantTypes[req.GrantType]
	if !ok {
		return nil, ErrGrantTypeUnSupported
	}

	return grant, nil
}

func (h *TokenHandler) getClient(req *TokenHandlerRequest) (*oauth2.Client, error) {
	client, err := h.Storage.GetClientWithSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		err = ErrInvalidClient.WithCause(err)
		h.Error("unable to get client",
			zap.String("client_id", client.ID),
			zap.String("client_secret", client.Secret),
			zap.Any("error", err),
		)

		return nil, err
	}

	if !client.IsGrantType(req.GrantType) {
		return nil, ErrGrantIsNotAllowed
	}

	return client, nil
}

func (h *TokenHandler) checkScope(grantName string, grantScopes, clientScopes []string) error {
	var err error

	switch grantName {
	case "client_credentials":
	case "authorization_code":
	case "refresh_token":
		return nil
	}

	ok, err := h.Storage.IsAvailableScope(grantScopes)
	if err != nil {
		err = ErrUnableCheckScope.WithCause(err)
		h.Error("unable get check scope from db",
			zap.Strings("scopes", grantScopes),
			zap.Any("error", err),
		)

		return err
	}
	if !ok {
		return ErrScopeUnSupported
	}

	for _, scope := range grantScopes {
		if !oauth2.HierarchicScope(scope, clientScopes) {
			return ErrScopeNotAllowed
		}
	}

	return nil
}

func (h *TokenHandler) createResponse(grantData *GrantData, accessToken, refreshToken string) TokenHandlerResponse {
	if accessToken == "" {
		return nil
	}

	data := TokenHandlerResponse{
		"access_token": accessToken,
		"token_type":   "bearer",
		"expires_in":   grantData.AccessTokenLifespan,
		"scope":        strings.Join(grantData.Scopes, " "),
	}

	if grantData.UserID != "" {
		data["user_id"] = grantData.UserID
	}

	if refreshToken != "" {
		data["refresh_token"] = refreshToken
	}

	for key, extra := range grantData.Extras {
		data[key] = extra
	}

	return data
}
