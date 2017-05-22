package token

import (
	"github.com/plimble/clover/oauth2"
)

var (
	ErrNotPostMethod            = oauth2.Error(400, "invalid_request", "The request method must be POST when requesting an access token")
	ErrParseForm                = oauth2.Error(400, "invalid_request", "unable to parse form")
	ErrGrantTypeRequired        = oauth2.Error(400, "invalid_request", "Missing parameters: grant_type required")
	ErrInvalidClient            = oauth2.Error(401, "unauthorized_client", "invalid client")
	ErrGrantTypeUnSupported     = oauth2.Error(400, "unsupported_grant_type", "grant type not supported")
	ErrGrantIsNotAllowed        = oauth2.Error(400, "invalid_grant", "grant type is not allowed")
	ErrScopeUnSupported         = oauth2.Error(400, "invalid_scope", "scope is not supported")
	ErrScopeNotAllowed          = oauth2.Error(400, "invalid_scope", "The scope requested is not allowed")
	ErrClientCredentialPublic   = oauth2.Error(400, "invalid_grant", "Public client is not allowed to use grant type client_credentials")
	ErrUnableCreateAccessToken  = oauth2.Error(500, "create_accesstoken", "Unable to create accesstoken")
	ErrUnableCreateRefreshToken = oauth2.Error(500, "create_refreshtoken", "Unable to create refreshtoken")
	ErrUnableCheckScope         = oauth2.Error(500, "check_scope", "Unable to check scope")
	ErrUnableGetAuthorizeCode   = oauth2.Error(500, "get_authorize", "Unable to get authorize code")
	ErrUnableGetRefreshToken    = oauth2.Error(500, "get_refreshtoken", "Unable to get refreshtoken")
	ErrUnableRevokeRefreshToken = oauth2.Error(500, "get_refreshtoken", "Unable to revoke refreshtoken")
	ErrCodeRequired             = oauth2.Error(400, "invalid_grant", "Missing parameter: code is required")
	ErrRedirectRequired         = oauth2.Error(400, "invalid_grant", "Missing parameter: redirect_uri is required")
	ErrClientIDMisMatched       = oauth2.Error(400, "invalid_grant", "client_id mismatched")
	ErrRedirectMisMatched       = oauth2.Error(400, "invalid_grant", "redirect_uri mismatched")
	ErrCodeExpired              = oauth2.Error(400, "invalid_grant", "code is expired")
	ErrRefreshTokenExpired      = oauth2.Error(400, "invalid_grant", "refreshtoken is expired")
	ErrUsernamePasswordRequired = oauth2.Error(400, "invalid_grant", "Missing parameters: username and password is required")
	ErrInvalidUser              = oauth2.Error(400, "invalid_grant", "username and password is invalid")
	ErrRefreshTokenRequired     = oauth2.Error(400, "invalid_grant", "Missing parameters: refresh_token required")
)
