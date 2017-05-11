package clover

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func redirectError(ctx *HTTPContext, redirectUrl string, err error) {
	u, _ := url.Parse(redirectUrl)
	query := u.Query()

	switch nerr := err.(type) {
	case *errorRes:
		query.Add("error", nerr.Code())
		query.Add("error", nerr.Message())
	default:
		query.Add("error", errUnKnownError.Code())
		query.Add("error_description", errUnKnownError.Message())
	}

	u.RawQuery = query.Encode()

	ctx.response.Header().Set("Cache-Control", "no-store")
	ctx.response.Header().Set("Pragma", "no-cache")

	http.Redirect(ctx.response, ctx.request, u.String(), http.StatusFound)
}

func redirect(ctx *HTTPContext, redirectUrl string) {
	ctx.response.Header().Set("Cache-Control", "no-store")
	ctx.response.Header().Set("Pragma", "no-cache")

	http.Redirect(ctx.response, ctx.request, redirectUrl, http.StatusFound)
}

func writeJson(ctx *HTTPContext, code int, v interface{}) {
	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.response.Header().Set("Cache-Control", "no-store")
	ctx.response.Header().Set("Pragma", "no-cache")
	ctx.response.WriteHeader(code)
	if v != nil {
		json.NewEncoder(ctx.response).Encode(v)
	}
}

func writeJsonError(ctx *HTTPContext, err error) error {
	if err == nil {
		return nil
	}

	switch nerr := err.(type) {
	case *errorRes:
		writeJson(ctx, nerr.HTTPCode(), nerr)
	default:
		writeJson(ctx, errUnKnownError.HTTPCode(), errUnKnownError)
	}

	return err
}
