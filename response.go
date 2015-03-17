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

type response struct {
	code        int
	data        map[string]interface{}
	inFragment  bool
	redirectURI string
	isErr       bool
	header      map[string]string
}

func (r *response) IsError() bool {
	return r.isErr
}

func (r *response) SetRedirect(uri, respType, state string) *response {
	r.redirectURI = uri
	r.code = 302

	if respType == RESP_TYPE_TOKEN {
		r.inFragment = true
	}

	if state != "" {
		r.data["state"] = state
	}

	return r
}

func (r *response) SetHeader(header map[string]string) {
	r.header = header
}

func (r *response) Write(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	for k, v := range r.header {
		w.Header().Set(k, v)
	}
	if r.code == 302 {
		r.redirect(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.code)
	encoder := json.NewEncoder(w)
	encoder.Encode(r.data)
}

func (r *response) redirect(w http.ResponseWriter) {
	u, err := url.Parse(r.redirectURI)
	if err != nil {
		errInternal(err.Error()).Write(w)
		return
	}

	q := u.Query()
	for key, desc := range r.data {
		q.Set(key, fmt.Sprint(desc))
	}
	if r.inFragment {
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

func NewRespErr(code int, key, desc string) *response {
	return &response{
		code:  code,
		isErr: true,
		data: map[string]interface{}{
			"error":             key,
			"error_description": desc,
		},
	}
}

func NewRespData(data map[string]interface{}) *response {
	return &response{
		code:  200,
		isErr: false,
		data:  data,
	}
}
