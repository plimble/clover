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
	logger  *zap.Logger
	Storage oauth2.Storage
}

func New(storage oauth2.Storage, logger *zap.Logger) *VerifyHandler {
	return &VerifyHandler{logger, storage}
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		oauth2.WriteJsonError(w, InvalidRequest("The request method must be POST"))
		return
	}

	if err := r.ParseForm(); err != nil {
		oauth2.WriteJsonError(w, InvalidRequest("Unable to parse form"))
		return
	}

	atoken := getAccessTokenFromRequest(r)

	req := &VerifyHandlerRequest{
		AccessToken: atoken,
		Scopes:      strings.Fields(r.Form.Get("scope")),
	}

	err := h.Verify(req)
	h.logger.Info("Verify AccessToken",
		zap.Any("VerifyHandlerRequest", req),
		zap.Any("error", err),
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
		if oauth2.IsNotFound(err) {
			return InvalidAccessToken("accesstoken is invalid or expired")
		}
		return err
	}

	if !at.Valid() {
		return InvalidAccessToken("accesstoken is expired")
	}

	if !at.HasScope(req.Scopes...) {
		return InvalidScope("scope request is not allowed")
	}

	return nil
}

func getAccessTokenFromRequest(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		return req.Form.Get("access_token")
	}

	return split[1]
}
