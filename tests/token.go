package tests

import (
	"encoding/json"
	"strings"
)

type AccessTokenResp struct {
	AccessToken  string                 `json:"access_token"`
	TokenType    string                 `json:"token_type"`
	ExpiresIn    int                    `json:"expires_in"`
	Scope        string                 `json:"scope"`
	Data         map[string]interface{} `json:"data"`
	RefreshToken string                 `json:"refresh_token"`
}

func getToken(body []byte) *AccessTokenResp {
	a := &AccessTokenResp{}
	if err := json.Unmarshal(body, a); err != nil {
		panic(err)
	}

	return a
}

func (a *AccessTokenResp) existScope(scope string) bool {
	scopes := strings.Fields(a.Scope)

	for i := 0; i < len(scopes); i++ {
		if scopes[i] == scope {
			return true
		}
	}

	return false
}
