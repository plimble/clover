package oauth2

import (
	"go.uber.org/zap/zapcore"
)

type AppErr struct {
	Message string `json:"error_description"`
	status  int
	Code    string `json:"error"`
	cause   error
}

func Error(message string, status int, code string, cause error) *AppErr {
	return &AppErr{message, status, code, cause}
}

func InvalidClient(message string) *AppErr {
	return &AppErr{message, 401, "invalid_client", nil}
}

func ServerError(message string, cause error) *AppErr {
	return &AppErr{message, 500, "server_error", cause}
}

func UnknownError(cause error) *AppErr {
	return &AppErr{"unknown error", 500, "server_error", cause}
}

func NotFound(err error) *AppErr {
	return &AppErr{"not found", 404, "not_found", err}
}

func (e *AppErr) Error() string { return e.Message }
func (e *AppErr) Cause() error  { return e.cause }
func (e *AppErr) WithCause(err error) *AppErr {
	e.cause = err
	return e
}

func (e *AppErr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("error_code", e.Code)
	enc.AddString("error_message", e.Message)
	if e.cause != nil {
		enc.AddString("error_cause", e.cause.Error())
	}

	return nil
}

func IsNotFound(err error) bool {
	switch t := err.(type) {
	case *AppErr:
		return t.status == 404
	}
	return false
}
