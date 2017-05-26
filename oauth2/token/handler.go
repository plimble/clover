package token

import (
	"net/http"
	"net/url"
	"strings"

	"fmt"

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
	grantTypes     map[string]GrantType
	Storage        oauth2.Storage
	TokenGenerator oauth2.TokenGenerator
}

func New(storage oauth2.Storage, tokenGen oauth2.TokenGenerator) *TokenHandler {
	return &TokenHandler{make(map[string]GrantType), storage, tokenGen}
}

func (h *TokenHandler) RegisterGrantType(grantType GrantType) {
	h.grantTypes[grantType.Name()] = grantType
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		oauth2.WriteJsonError(w, InvalidRequest("The request method must be POST"))
		return
	}

	auth := r.Header.Get("Authorization")
	clientID, clientSecret, err := oauth2.GetCredentialsFromHttp(auth)
	if err != nil {
		oauth2.WriteJsonError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		oauth2.WriteJsonError(w, InvalidRequest("unable to parse form"))
		return
	}

	req := &TokenHandlerRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    r.Form.Get("grant_type"),
		Form:         r.Form,
	}

	res, err := h.RequestToken(req)
	oauth2.Logger.Info("RequestToken",
		zap.Any("TokenHandlerRequest", req),
		zap.Any("TokenHandlerResponse", res),
		zap.Any("error", err),
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
		return nil, InvalidRequest("Missing parameters: grant_type required")
	}

	grant, ok := h.grantTypes[req.GrantType]
	if !ok {
		return nil, UnsupportedGrantType(fmt.Sprintf("server is not supported %s grant type", req.GrantType))
	}

	return grant, nil
}

func (h *TokenHandler) getClient(req *TokenHandlerRequest) (*oauth2.Client, error) {
	client, err := h.Storage.GetClientWithSecret(req.ClientID, req.ClientSecret)
	if err != nil {
		if oauth2.IsNotFound(err) {
			return nil, UnauthorizedClient("The authenticated client is not authorized")
		}

		return nil, err
	}

	if !client.HasGrantType(req.GrantType) {
		return nil, UnsupportedGrantType(fmt.Sprintf("client is not supported %s grant type", req.GrantType))
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
		return err
	}
	if !ok {
		return InvalidScope("server is not supported your scope")
	}

	for _, scope := range grantScopes {
		if !oauth2.HierarchicScope(scope, clientScopes) {
			return InvalidScope("scope is not allowed")
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
