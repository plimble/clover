package verify

import (
	"github.com/plimble/clover/oauth2"
)

var (
	ErrNotPostMethod        = oauth2.Error(400, "invalid_request", "The request method must be POST when requesting an access token")
	ErrParseForm            = oauth2.Error(400, "invalid_request", "unable to parse form")
	ErrUnableGetAccessToken = oauth2.Error(500, "get_accesstoken", "unable to get accesstoken")
	ErrInvalidAccessToken   = oauth2.Error(500, "invalid_accesstoken", "invalid accesstoken")
	ErrInvalidScope         = oauth2.Error(500, "invalid_scope", "invalid scope")
)
