package oauth2

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func WriteJsonError(w http.ResponseWriter, err error) error {
	switch nerr := err.(type) {
	case *AppErr:
		WriteJson(w, nerr.status, nerr)
	default:
		WriteJson(w, 500, UnknownError())
	}

	return err
}
