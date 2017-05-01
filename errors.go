package clover

import (
	"fmt"
	"io"

	"github.com/plimble/errors"
)

func ErrGrantTypeNotFound() error {
	return NewError(400, "invalid_request", "The request method must be POST when requesting an access token")
}

func ErrGrantTypeNotSupport(grantType string) error {
	return NewErrorf(400, "unsupported_grant_type", `Grant type "%s" not supported`, grantType)
}

func ErrInvalidClient(id string, err error) error {
	return NewErrorfWithCause(err, 400, "invalid_grant", "%s doesn not exist or is invalid for the client", id)
}

func ErrClientCredentialRequired() error {
	return NewError(400, "invalid_client", "client credentials are required")
}

func ErrInternalServer(cause error) error {
	return NewErrorWithCause(cause, 500, "internal_server", "Internal Server Error")
}

func ErrInvalidAuthMessage() error {
	return NewError(400, "invalid_client", "Invalid authorization message")
}

func ErrInvalidAuthHeader() error {
	return NewError(400, "invalid_client", "Invalid authorization header")
}

func ErrInvalidParseForm(cause error) error {
	return NewErrorWithCause(cause, 400, "invalid_request", "Invalid parse form")
}

func ErrUnSupportedScope() error {
	return NewError(400, "invalid_scope", "An unsupported scope was requested")
}

func ErrInvalidScope() error {
	return NewError(400, "invalid_scope", "The scope requested is invalid for this request")
}

type Error struct {
	cause         error
	ErrMessage    string `json:"message"`
	ErrCode       string `json:"error"`
	ErrHttpCode   int    `json:"-"`
	*errors.Stack `json:"-"`
}

func (e Error) Error() string   { return e.ErrMessage }
func (e Error) HTTPCode() int   { return e.ErrHttpCode }
func (e Error) Message() string { return e.ErrMessage }
func (e Error) Code() string    { return e.ErrCode }
func (e Error) Cause() error    { return e.cause }

func NewError(httpstatus int, code, message string) error {
	return Error{nil, message, code, httpstatus, errors.Callers()}
}

func NewErrorf(httpstatus int, code, message string, v ...interface{}) error {
	return Error{nil, fmt.Sprintf(message, v...), code, httpstatus, errors.Callers()}
}

func NewErrorWithCause(cause error, httpstatus int, code, message string) error {
	return Error{cause, message, code, httpstatus, errors.Callers()}
}

func NewErrorfWithCause(cause error, httpstatus int, code, message string, v ...interface{}) error {
	return Error{cause, fmt.Sprintf(message, v...), code, httpstatus, errors.Callers()}
}

func (e Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.ErrMessage)
			fmt.Fprintf(s, "%+v", e.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.ErrMessage)
	}
}
