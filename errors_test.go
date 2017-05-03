package clover

import (
	"testing"

	basicerr "errors"
	"fmt"

	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err := aa()
	fmt.Printf("%+v\n", err)

	xerr := errors.WithStack(ErrInvalidRefreshToken)
	fmt.Printf("%+v\n", xerr)
}

func aa() error {
	err := bb()
	return ErrGrantTypeNotFound.WithCause(err)
}

func bb() error {
	err := cc()
	return errors.WithStack(err)
}

func cc() error {
	return basicerr.New("xxxx")
}
