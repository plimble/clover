package oauth2

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type AppErr struct {
	Message string `json:"message"`
	status  int
	Code    string `json:"error"`
	cause   error
}

func Error(status int, errCode, msg string) *AppErr {
	return &AppErr{msg, status, errCode, nil}
}

func Errorf(status int, errCode, format string, v ...interface{}) *AppErr {
	return Error(status, errCode, fmt.Sprintf(format, v...))
}

func UnknownError() *AppErr {
	return &AppErr{"unknown error", 500, "unknown", nil}
}

func DbNotFoundError(err error) *AppErr {
	return &AppErr{"db not found", 404, "not_found", err}
}

func (e *AppErr) Error() string { return e.Message }
func (e *AppErr) Cause() error  { return e.cause }
func (e *AppErr) WithCause(err error) *AppErr {
	return &AppErr{e.Message, e.status, e.Code, err}
}

func (e *AppErr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("error_code", e.Code)
	enc.AddString("error_message", e.Message)
	if cause := ErrorCause(e); cause != nil {
		enc.AddString("error_cause", cause.Error())
	}

	return nil
}

func ErrorCause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		c, ok := err.(causer)
		if !ok {
			break
		}
		cause := c.Cause()
		if cause == nil {
			break
		}
		err = cause
	}
	return err
}
