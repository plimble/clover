package clover

import (
	"testing"
)

func TestStackError(t *testing.T) {
	err := NewError(500, "error_code", "error!!")
	nerr, ok := err.(Error)
	if !ok {
		t.Error("error should be Error struct")
	}

	if nerr.Code() != "error_code" {
		t.Error(`error code should be "error_code"`)
	}

	if nerr.HTTPCode() != 500 {
		t.Error(`error http code should be 500`)
	}

	if nerr.Message() != "error!!" {
		t.Error(`error http code should be "error!!"`)
	}
}
