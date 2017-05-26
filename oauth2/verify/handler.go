package verify

import (
	"net/http"
	"strings"

	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type VerifyHandlerRequest struct {
	AccessToken string
	Scopes      []string
}

type VerifyHandler struct {
	*zap.Logger
	Storage oauth2.Storage
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		oauth2.WriteJsonError(w, ErrNotPostMethod)
		return
	}

	atoken := getAccessTokenFromRequest(r)

	if err := r.ParseForm(); err != nil {
		oauth2.WriteJsonError(w, ErrParseForm)
		return
	}

	req := &VerifyHandlerRequest{
		AccessToken: atoken,
		Scopes:      strings.Fields(r.Form.Get("scope")),
	}

	err := h.Verify(req)
	h.Logger.Info("Verify",
		zap.Any("VerifyHandlerRequest", req),
		zap.Error(err),
	)
	if err != nil {
		oauth2.WriteJsonError(w, err)
		return
	}

	oauth2.WriteJson(w, 200, nil)
}

func (h *VerifyHandler) Verify(req *VerifyHandlerRequest) error {
	at, err := h.Storage.GetAccessToken(req.AccessToken)
	if err != nil {
		return ErrUnableGetAccessToken.WithCause(err)
	}

	if !at.Valid() {
		return ErrInvalidAccessToken
	}

	if !at.HasScope(req.Scopes...) {
		return ErrInvalidScope
	}

	return nil
}

func getAccessTokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		err := req.ParseForm()
		if err != nil {
			return ""
		}
		return req.Form.Get("access_token")
	}

	return split[1]
}
