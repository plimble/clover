package clover

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setUp() *Clover {
	return New(&Config{})
}

func TestValidateAuthorize(t *testing.T) {
	c := setUp()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ar := c.ValidateAuthorize(w, r)

	fmt.Println(ar)
	fmt.Println(w.Body.String())
	assert.NotNil(t, ar)
}
