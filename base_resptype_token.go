package clover

import (
	"strings"
)

type createTokenFunc func(clientID, userID string, scopes []string, data map[string]interface{}) (*AccessToken, *Response)

type baseTokenRespType struct {
	accessTokenstore  AccessTokenStore
	refreshTokenStore RefreshTokenStore
	accessLifeTime    int
	refreshLifeTime   int
}

func newBaseTokenRespType(accessTokenstore AccessTokenStore, refreshTokenStore RefreshTokenStore, accessLifeTime, refreshLifeTime int) *baseTokenRespType {
	return &baseTokenRespType{accessTokenstore, refreshTokenStore, accessLifeTime, refreshLifeTime}
}

func (rt *baseTokenRespType) createAccessToken(clientID, userID string, scopes []string, data map[string]interface{}, token string) (*AccessToken, *Response) {
	aToken := &AccessToken{
		AccessToken: token,
		ClientID:    clientID,
		UserID:      userID,
		Expires:     addSecondUnix(rt.accessLifeTime),
		Scope:       scopes,
		Data:        data,
	}

	if err := rt.accessTokenstore.SetAccessToken(aToken); err != nil {
		return nil, errInternal(err.Error())
	}

	return aToken, nil
}

func (rt *baseTokenRespType) createRefreshToken(clientID, userID string, scopes []string, data map[string]interface{}, refreshToken string, includeRefresh bool) (string, *Response) {
	if !includeRefresh || rt.refreshTokenStore == nil || rt.refreshLifeTime < 1 {
		return "", nil
	}

	r := &RefreshToken{
		RefreshToken: refreshToken,
		ClientID:     clientID,
		UserID:       userID,
		Expires:      addSecondUnix(rt.refreshLifeTime),
		Scope:        scopes,
		Data:         data,
	}

	if err := rt.refreshTokenStore.SetRefreshToken(r); err != nil {
		return "", errInternal(err.Error())
	}

	return r.RefreshToken, nil
}

func (rt *baseTokenRespType) createRespData(token string, scopes []string, refresh, state string, data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   rt.accessLifeTime,
		"scope":        strings.Join(scopes, " "),
	}

	if len(data) > 0 {
		result["data"] = data
	}

	if refresh != "" {
		result["refresh_token"] = refresh
	}

	if state != "" {
		result["state"] = state
	}

	return result
}
