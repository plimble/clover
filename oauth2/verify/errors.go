package verify

import "github.com/plimble/clover/oauth2"

// InvalidRequest The request is missing a required parameter, includes an
// unsupported parameter value
func InvalidRequest(message string) *oauth2.AppErr {
	return oauth2.Error(message, 400, "invalid_request", nil)
}

// InvalidAccessToken The request send an invalid access token or expired
func InvalidAccessToken(message string) *oauth2.AppErr {
	return oauth2.Error(message, 400, "invalid_accesstoken", nil)
}

// InvalidScope The requested scope is not allowed
func InvalidScope(message string) *oauth2.AppErr {
	return oauth2.Error(message, 401, "invalid_scope", nil)
}
