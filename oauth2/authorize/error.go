package authorize

import "github.com/plimble/clover/oauth2"

// InvalidRequest The request is missing a required parameter, includes an
// invalid parameter value, includes a parameter more than
// once, or is otherwise malformed.
func InvalidRequest(message string) *oauth2.AppErr {
	return oauth2.Error(message, 400, "invalid_request", nil)
}

// UnauthorizedClient The client is not authorized to request an authorization
// code using this method.
func UnauthorizedClient(message string) *oauth2.AppErr {
	return oauth2.Error(message, 401, "unauthorized_client", nil)
}

// AccessDenied The resource owner or authorization server denied the request.
func AccessDenied(message string) *oauth2.AppErr {
	return oauth2.Error(message, 403, "access_denied", nil)
}

// UnsupportedResponseType The authorization server does not support obtaining an
// authorization code using this method.
func UnsupportedResponseType(message string) *oauth2.AppErr {
	return oauth2.Error(message, 400, "unsupported_response_type", nil)
}

// InvalidScope The requested scope is invalid, unknown, or malformed.
func InvalidScope(message string) *oauth2.AppErr {
	return oauth2.Error(message, 400, "invalid_scope", nil)
}
