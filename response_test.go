package clover

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestResponse_Write(t *testing.T) {
	w := httptest.NewRecorder()

	resp := newRespData(map[string]interface{}{
		"name": "Tester",
	})

	resp.Write(w)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "{\"name\":\"Tester\"}\n", w.Body.String())
}

func TestResponse_Write_Redirect(t *testing.T) {
	w := httptest.NewRecorder()

	resp := newRespData(map[string]interface{}{
		"name": "Tester",
	})

	resp.setRedirect("http://localhost/redirect", RESP_TYPE_CODE, "1234")

	resp.Write(w)
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "http://localhost/redirect?name=Tester&state=1234", w.Header().Get("Location"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestResponse_Write_RedirectFragment(t *testing.T) {
	w := httptest.NewRecorder()

	resp := newRespData(map[string]interface{}{
		"name": "Tester",
	})

	resp.setRedirect("http://localhost/redirect", RESP_TYPE_TOKEN, "1234")

	resp.Write(w)
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "http://localhost/redirect#name=Tester&state=1234", w.Header().Get("Location"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestResponse_SetRedirect_Fragment(t *testing.T) {
	resp := newRespData(map[string]interface{}{
		"name": "Tester",
	})

	resp.setRedirect("http://localhost/redirect", RESP_TYPE_TOKEN, "1234")

	assert.Equal(t, 302, resp.code)
	assert.True(t, resp.isFragment)
	assert.Equal(t, "1234", resp.data["state"])
}

func TestResponse_SetRedirect(t *testing.T) {
	resp := newRespData(map[string]interface{}{
		"name": "Tester",
	})

	resp.setRedirect("http://localhost/redirect", RESP_TYPE_CODE, "1234")

	assert.Equal(t, 302, resp.code)
	assert.False(t, resp.isFragment)
	assert.Equal(t, "1234", resp.data["state"])
}
