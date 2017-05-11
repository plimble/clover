package clover

import (
	"fmt"
	"io"
)

var (
	errUnKnownError             = Error(500, "unknown", "Unknown error")
	errInvalidTokenFormat       = Error(400, "invalid_token", "Invalid token")
	errTokenSignatureMismatch   = Error(400, "mismatch_signature", "Signature mismatch")
	errInvalidParseForm         = Error(400, "invalid_request", "Invalid parse form")
	errInvalidScope             = Error(400, "invalid_scope", "The scope requested is invalid for this request")
	errInvalidAuthHeader        = Error(400, "invalid_client", "Invalid authorization header")
	errInvalidAuthMessage       = Error(400, "invalid_client", "Invalid authorization message")
	errClientCredentialRequired = Error(400, "invalid_client", "client credentials are required")
	errGrantTypeNotFound        = Error(400, "invalid_request", "The request method must be POST when requesting an access token")
	errInternalServer           = Error(500, "internal_server", "Internal Server Error")
	errPublicClient             = Error(400, "invalid_client", "The client is public and thus not allowed to use grant type client_credentials")
	errInvalidUsernamePassword  = Error(400, "invalid_grant", "Invalid username or password")
	errUsernamePasswordRequired = Error(400, "invalid_grant", `Missing parameters: "username" and "password" required`)
	errInvalidAccessToken       = Error(400, "invalid_token", `The access token provided is invalid`)
	errInvalidRefreshToken      = Error(400, "invalid_grant", `The refresh token provided is invalid`)
	errInvalidAuthCode          = Error(400, "invalid_grant", `Authorization code does not exist or is invalid for the client`)
	errRefreshTokenRequired     = Error(400, "invalid_request", `Missing parameter: "refresh_token" is required`)
	errRefreshTokenExpired      = Error(400, "invalid_grant", `Refresh token has expired`)
	errCodeRequired             = Error(400, "invalid_request", `Missing parameter: "code" is required`)
	errRedirectURIRequired      = Error(400, "invalid_request", `Missing parameter: "redirect_uri" is required`)
	errRedirectMisMatch         = Error(400, "redirect_uri_mismatch", `The redirect URI is missing or do not match`)
	errCodeExpired              = Error(400, "invalid_grant", `The authorization code has expire`)
	errClientIDMisMatch         = Error(400, "invalid_request", `Client mismatched`)
	errRevokeTokenTypeInvalid   = Error(400, "invalid_request", `Token type hint must be either "access_token" or "refresh_token"`)
	errRevokeTokenRequired      = Error(400, "invalid_request", `Missing token parameter to revoke`)
	errMethodPostRequired       = Error(405, "invalid_request", `The request method must be POST when requesting an access token`)
	errIntroTokenRequired       = Error(400, "invalid_request", `Missing token parameter to introspection`)
	errIntroTokenTypeRequired   = Error(400, "invalid_request", `Missing token type parameter to introspection`)
	errIntroTokenTypeInvalid    = Error(400, "invalid_request", `Token type hint must be either "access_token" or "refresh_token"`)
	errNoClient                 = Error(400, "invalid_client", "No client id supplied")
	errResponseTypeRequired     = Error(400, "invalid_request", "response_type is required")
	errStateRequired            = Error(400, "invalid_request", "state is required")
	errInvalidSession           = Error(400, "invalid_request", "Invalid session")
	errInvalidChallenge         = Error(400, "invalid_request", "Invalid challenge")
	errUserIDRequired           = Error(400, "invalid_request", "Invalid user id is required")
	errChallengeExpired         = Error(400, "invalid_request", `challenge token has expire`)
	errResponseTypeUnSupported  = Error(400, "invalid_request", `Response type unsupported`)
	errInvalidClient            = Error(400, "invalid_grant", `Invalid client`)
)

func ErrGrantTypeNotSupport(grantType string) *errorRes {
	return Errorf(400, "unsupported_grant_type", `Grant type "%s" not supported`, grantType)
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
