package clover

import (
	"fmt"
	"io"
)

var (
	ErrInvalidTokenFormat       = Error(400, "invalid_token", "Invalid token")
	ErrTokenSignatureMismatch   = Error(400, "mismatch_signature", "Signature mismatch")
	ErrInvalidParseForm         = Error(400, "invalid_request", "Invalid parse form")
	ErrUnSupportedScope         = Error(400, "invalid_scope", "An unsupported scope was requested")
	ErrInvalidScope             = Error(400, "invalid_scope", "The scope requested is invalid for this request")
	ErrInvalidAuthHeader        = Error(400, "invalid_client", "Invalid authorization header")
	ErrInvalidAuthMessage       = Error(400, "invalid_client", "Invalid authorization message")
	ErrClientCredentialRequired = Error(400, "invalid_client", "client credentials are required")
	ErrNotImplemented           = Error(501, "not_implemented", "Not Implemented")
	ErrGrantTypeNotFound        = Error(400, "invalid_request", "The request method must be POST when requesting an access token")
	ErrInternalServer           = Error(500, "internal_server", "Internal Server Error")
	ErrSecretNotStrong          = Error(500, "invalid_secret", "Secret is not strong enough")
	ErrPublicClient             = Error(400, "invalid_client", "The client is public and thus not allowed to use grant type client_credentials")
	ErrInvalidUsernamePassword  = Error(400, "invalid_grant", "Invalid username or password")
	ErrUsernamePasswordRequired = Error(400, "invalid_grant", `Missing parameters: "username" and "password" required`)
	ErrAccessTokenExpired       = Error(400, "invalid_token", `Access token is expired`)
	ErrInvalidAccessToken       = Error(400, "invalid_token", `The access token provided is invalid`)
	ErrInvalidRefreshToken      = Error(400, "invalid_grant", `The refresh token provided is invalid`)
	ErrInvalidAuthCode          = Error(400, "invalid_grant", `Authorization code does not exist or is invalid for the client`)
	ErrRefreshTokenRequired     = Error(400, "invalid_request", `Missing parameter: "refresh_token" is required`)
	ErrRefreshTokenExpired      = Error(400, "invalid_grant", `Refresh token has expired`)
	ErrCodeRequired             = Error(400, "invalid_request", `Missing parameter: "code" is required`)
	ErrRedirectMisMatch         = Error(400, "redirect_uri_mismatch", `The redirect URI is missing or do not match`)
	ErrCodeExpired              = Error(400, "invalid_grant", `The authorization code has expire`)
	ErrClientIDMisMatch         = Error(400, "invalid_request", `The request method must be POST when revoking an access token`)
	ErrRevokeTokenTypeInvalid   = Error(400, "invalid_request", `Token type hint must be either "access_token" or "refresh_token"`)
	ErrRevokeTokenRequired      = Error(400, "invalid_request", `Missing token parameter to revoke`)
)

func ErrGrantTypeNotSupport(grantType string) *errorRes {
	return Errorf(400, "unsupported_grant_type", `Grant type "%s" not supported`, grantType)
}

func ErrInvalidClient(id string) *errorRes {
	return Errorf(400, "invalid_grant", "%s doesn not exist or is invalid for the client", id)
}

type errorRes struct {
	cause       error
	ErrMessage  string `json:"message"`
	ErrCode     string `json:"error"`
	ErrHttpCode int    `json:"-"`
}

func (e *errorRes) Error() string   { return e.ErrMessage }
func (e *errorRes) HTTPCode() int   { return e.ErrHttpCode }
func (e *errorRes) Message() string { return e.ErrMessage }
func (e *errorRes) Code() string    { return e.ErrCode }
func (e *errorRes) Cause() error    { return e.cause }
func (e *errorRes) WithCause(err error) *errorRes {
	e.cause = err
	return e
}

func (e *errorRes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if e.Cause() != nil {
				fmt.Fprintln(s, e.ErrMessage)
				fmt.Fprintf(s, "%+v\n", e.Cause())
			} else {
				fmt.Fprint(s, e.ErrMessage)
			}
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}

func Error(httpstatus int, code, message string) *errorRes {
	return &errorRes{nil, message, code, httpstatus}
}

func Errorf(httpstatus int, code, message string, v ...interface{}) *errorRes {
	return &errorRes{nil, fmt.Sprintf(message, v...), code, httpstatus}
}
