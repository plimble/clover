package clover

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	CODE_BAD_REQUEST       = 400
	CODE_INTERNAL          = 500
	CODE_UNAUTH            = 401
	CODE_PERMISSION_DENIED = 403
	CODE_REDIRECT          = 302
)

type Response struct {
	code        int
	isFragment  bool
	redirectURI string
	isErr       bool
	data        map[string]interface{}
	state       string
}

func (r *Response) IsError() bool {
	return r.isErr
}

func (r *Response) IsRedirect() bool {
	return r.code == 302
}

func (r *Response) setRedirect(uri string, fragment bool, state string) *Response {
	r.redirectURI = uri
	r.code = 302

	r.isFragment = fragment
	if state != "" {
		r.data["state"] = state
	}

	return r
}

func (r *Response) Write(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	if r.code == 302 {
		r.redirect(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.code)
	encoder := json.NewEncoder(w)
	encoder.Encode(r.data)
}

func (r *Response) redirect(w http.ResponseWriter) {
	u, err := url.Parse(r.redirectURI)
	if err != nil {
		errInternal(err.Error()).Write(w)
		return
	}

	q := u.Query()
	for key, desc := range r.data {
		q.Set(key, fmt.Sprint(desc))
	}
	if r.isFragment {
		u.RawQuery = ""
		u.Fragment, err = url.QueryUnescape(q.Encode())
		if err != nil {
			errInternal(err.Error()).Write(w)
			return
		}
	} else {
		u.RawQuery = q.Encode()
	}

	w.Header().Set("Location", u.String())
	w.WriteHeader(302)
}

func NewErrResp(code int, title, desc string) *Response {
	return &Response{
		code:  code,
		isErr: true,
		data: map[string]interface{}{
			"error":             title,
			"error_description": desc,
		},
	}
}

func NewRespData(data map[string]interface{}) *Response {
	return &Response{
		code:  200,
		isErr: false,
		data:  data,
	}
}
