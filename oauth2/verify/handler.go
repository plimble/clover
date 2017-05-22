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

	atoken := oauth2.GetAccessTokenFromRequest(r)

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

	for _, scope := range req.Scopes {
		if !oauth2.HierarchicScope(scope, at.Scopes) {
			return ErrInvalidScope
		}
	}

	return nil
}
