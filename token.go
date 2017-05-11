package clover

import (
	"strings"
	"sync"
)

// TokenManager Manage acess, refresh, authroize token
//go:generate mockery -name TokenManager
type TokenManager interface {
	GenerateAccessToken(ctx *AccessTokenContext, includeRefreshToken bool) (*AccessToken, *RefreshToken, error)
	GenerateRefreshToken(ctx *AccessTokenContext) (*RefreshToken, error)
	GenerateAuthorizeCode(ctx *AuthorizeContext, authorizeCodeLifespan int) (*AuthorizeCode, error)
	GetAccessToken(token string) (*AccessToken, error)
	GetRefreshToken(token string) (*RefreshToken, error)
	GetAuthorizeCode(code string) (*AuthorizeCode, error)
	DeleteAccessToken(token string) error
	DeleteRefreshToken(token string) error
	DeleteAuthorizeCode(code string) error
}

//go:generate mockery -name TokenStorage
type TokenStorage interface {
	DeleteAccessToken(token string) error
	SaveAccessToken(accessToken *AccessToken) error
	GetAccessToken(token string) (*AccessToken, error)

	DeleteRefreshToken(token string) error
	SaveRefreshToken(refreshToken *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)

	DeleteAuthorizeCode(code string) error
	SaveAuthorizeCode(authCode *AuthorizeCode) error
	GetAuthorizeCode(code string) (*AuthorizeCode, error)
}

//go:generate mockery -name TokenGenerator
type TokenGenerator interface {
	Generate(clientID, userID, scope string, tokenLifespan int) (string, error)
}

type AccessToken struct {
	AccessToken string   `json:"access_token"`
	ClientID    string   `json:"client_id,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	Expired     int64    `json:"expired"`
	Scopes      []string `json:"scopes"`
}

type AuthorizeCode struct {
	Code         string   `json:"code"`
	ClientID     string   `json:"client_id"`
	UserID       string   `json:"user_id"`
	Expired      int64    `json:"expired"`
	Scopes       []string `json:"scopes"`
	RedirectURI  string   `json:"redirect_uri"`
	ResponseType string   `json:"response_type"`
}

type RefreshToken struct {
	RefreshToken         string   `json:"refresh_token"`
	ClientID             string   `json:"client_id"`
	UserID               string   `json:"user_id"`
	Expired              int64    `json:"expired"`
	Scopes               []string `json:"scopes"`
	AccessTokenLifespan  int      `json:"access_token_lifespan"`
	RefreshTokenLifespan int      `json:"refresh_token_lifespan"`
}

type tokenManager struct {
	accessTokenGenerator  TokenGenerator
	refreshTokenGenerator TokenGenerator
	authCodeGenerator     TokenGenerator
	tokenStore            TokenStorage
}

func NewTokenManager(accessTokenGenerator, refreshTokenGenerator, authCodeGenerator TokenGenerator, tokenStore TokenStorage) TokenManager {
	return &tokenManager{accessTokenGenerator, refreshTokenGenerator, authCodeGenerator, tokenStore}
}

func (t *tokenManager) GenerateAccessToken(ctx *AccessTokenContext, includeRefreshToken bool) (*AccessToken, *RefreshToken, error) {
	var wg sync.WaitGroup
	var accessToken *AccessToken
	var accessTokenErr error
	var acToken string

	wg.Add(1)
	go func() {
		defer wg.Done()
		acToken, accessTokenErr = t.accessTokenGenerator.Generate(
			ctx.Client.ID,
			ctx.UserID,
			strings.Join(ctx.Scopes, " "),
			ctx.AccessTokenLifespan,
		)
		if accessTokenErr != nil {
			return
		}

		accessToken = &AccessToken{
			AccessToken: acToken,
			ClientID:    ctx.Client.ID,
			UserID:      ctx.UserID,
			Expired:     addSecondUnix(ctx.AccessTokenLifespan),
			Scopes:      ctx.Scopes,
		}

		accessTokenErr = t.tokenStore.SaveAccessToken(accessToken)
	}()

	var refreshTokenErr error
	var refreshToken *RefreshToken

	if includeRefreshToken {
		wg.Add(1)
		go func() {
			refreshToken, refreshTokenErr = t.GenerateRefreshToken(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
	if accessTokenErr != nil {
		return nil, nil, accessTokenErr
	}
	if refreshTokenErr != nil {
		return nil, nil, refreshTokenErr
	}

	return accessToken, refreshToken, nil
}

func (t *tokenManager) GenerateRefreshToken(ctx *AccessTokenContext) (*RefreshToken, error) {
	token, err := t.refreshTokenGenerator.Generate(
		ctx.Client.ID,
		ctx.UserID,
		strings.Join(ctx.Scopes, " "),
		ctx.RefreshTokenLifespan,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := &RefreshToken{
		RefreshToken:         token,
		ClientID:             ctx.Client.ID,
		UserID:               ctx.UserID,
		Expired:              addSecondUnix(ctx.RefreshTokenLifespan),
		Scopes:               ctx.Scopes,
		AccessTokenLifespan:  ctx.AccessTokenLifespan,
		RefreshTokenLifespan: ctx.RefreshTokenLifespan,
	}

	err = t.tokenStore.SaveRefreshToken(refreshToken)

	return refreshToken, err
}

func (t *tokenManager) GenerateAuthorizeCode(ctx *AuthorizeContext, authorizeCodeLifespan int) (*AuthorizeCode, error) {
	token, err := t.authCodeGenerator.Generate(
		ctx.Client.ID,
		ctx.UserID,
		strings.Join(ctx.Scopes, " "),
		authorizeCodeLifespan,
	)
	if err != nil {
		return nil, err
	}

	authCode := &AuthorizeCode{
		Code:         token,
		ClientID:     ctx.Client.ID,
		UserID:       ctx.UserID,
		Expired:      addSecondUnix(authorizeCodeLifespan),
		Scopes:       ctx.Scopes,
		RedirectURI:  ctx.RedirectURI,
		ResponseType: ctx.ResponseType,
	}

	err = t.tokenStore.SaveAuthorizeCode(authCode)

	return authCode, err
}

func (t *tokenManager) GetAccessToken(token string) (*AccessToken, error) {
	return t.tokenStore.GetAccessToken(token)
}

func (t *tokenManager) GetRefreshToken(token string) (*RefreshToken, error) {
	return t.tokenStore.GetRefreshToken(token)
}

func (t *tokenManager) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	return t.tokenStore.GetAuthorizeCode(code)
}

func (t *tokenManager) DeleteAccessToken(token string) error {
	return t.tokenStore.DeleteAccessToken(token)
}

func (t *tokenManager) DeleteRefreshToken(token string) error {
	return t.tokenStore.DeleteRefreshToken(token)
}

func (t *tokenManager) DeleteAuthorizeCode(code string) error {
	return t.tokenStore.DeleteAuthorizeCode(code)
}
