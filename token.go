package clover

import "github.com/pkg/errors"

//go:generate mockery -name TokenGenerator -case underscore
type TokenGenerator interface {
	Generate() (string, error)
	Validate(token string) error
}

//go:generate mockery -name TokenManager -case underscore
type TokenManager interface {
	GenerateAccessToken(clientID, userID string, expiresIn int, scopes []string) (*AccessToken, error)
	GenerateRefreshToken(clientID, userID string, expiresIn int, scopes []string) (*RefreshToken, error)
	GenerateAuthorizeCode(clientID, userID string, expiresIn int, scopes []string, redirectUri string) (*AuthorizeCode, error)
	GetAccessToken(token string) (*AccessToken, error)
	GetRefreshToken(token string) (*RefreshToken, error)
	GetAuthorizeCode(code string) (*AuthorizeCode, error)
	DeleteAccessToken(token string) error
	DeleteRefreshToken(token string) error
	DeleteAuthorizeCode(code string) error
}

type tokenManager struct {
	generator  TokenGenerator
	tokenStore TokenStore
}

func NewTokenManager(generator TokenGenerator, tokenStore TokenStore) TokenManager {
	return &tokenManager{generator, tokenStore}
}

func (t *tokenManager) GenerateAccessToken(clientID, userID string, expiresIn int, scopes []string) (*AccessToken, error) {
	token, err := t.generator.Generate()
	if err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	accessToken := &AccessToken{
		AccessToken: token,
		ClientID:    clientID,
		UserID:      userID,
		Expired:     addSecondUnix(expiresIn),
		Scopes:      scopes,
	}

	if err = t.tokenStore.SaveAccessToken(accessToken); err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return accessToken, nil
}

func (t *tokenManager) GenerateRefreshToken(clientID, userID string, expiresIn int, scopes []string) (*RefreshToken, error) {
	token, err := t.generator.Generate()
	if err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	refreshToken := &RefreshToken{
		RefreshToken: token,
		ClientID:     clientID,
		UserID:       userID,
		Expired:      addSecondUnix(expiresIn),
		Scopes:       scopes,
	}

	if err = t.tokenStore.SaveRefreshToken(refreshToken); err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return refreshToken, nil
}

func (t *tokenManager) GenerateAuthorizeCode(clientID, userID string, expiresIn int, scopes []string, redirectUri string) (*AuthorizeCode, error) {
	token, err := t.generator.Generate()
	if err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	authCode := &AuthorizeCode{
		Code:        token,
		ClientID:    clientID,
		UserID:      userID,
		Expired:     addSecondUnix(expiresIn),
		Scopes:      scopes,
		RedirectURI: redirectUri,
	}

	if err = t.tokenStore.SaveAuthCode(authCode); err != nil {
		return nil, ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return authCode, nil
}

func (t *tokenManager) GetAccessToken(token string) (*AccessToken, error) {
	accessToken, err := t.tokenStore.GetAccessToken(token)
	if err != nil {
		return nil, ErrInvalidAccessToken.WithCause(errors.WithStack(err))
	}

	return accessToken, nil
}

func (t *tokenManager) GetRefreshToken(token string) (*RefreshToken, error) {
	refreshToken, err := t.tokenStore.GetRefreshToken(token)
	if err != nil {
		return nil, ErrInvalidRefreshToken.WithCause(errors.WithStack(err))
	}

	return refreshToken, nil
}

func (t *tokenManager) GetAuthorizeCode(code string) (*AuthorizeCode, error) {
	authCode, err := t.tokenStore.GetAuthCode(code)
	if err != nil {
		return nil, ErrInvalidAuthCode.WithCause(errors.WithStack(err))
	}

	return authCode, nil
}

func (t *tokenManager) DeleteAccessToken(token string) error {
	err := t.tokenStore.DeleteAccessToken(token)
	if err != nil {
		return ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return nil
}

func (t *tokenManager) DeleteRefreshToken(token string) error {
	err := t.tokenStore.DeleteRefreshToken(token)
	if err != nil {
		return ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return nil
}

func (t *tokenManager) DeleteAuthorizeCode(code string) error {
	err := t.tokenStore.DeleteAuthCode(code)
	if err != nil {
		return ErrInternalServer.WithCause(errors.WithStack(err))
	}

	return nil
}
