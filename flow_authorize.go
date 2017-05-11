package clover

import (
	"net/url"

	"strings"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

type AuthorizeRes struct {
	Error       error
	RedirectURI *url.URL
}

type authorizeFlow struct {
	resTypes       map[string]ResponseType
	tokenManager   TokenManager
	clientStorage  ClientStorage
	consent        Consent
	session        sessions.Store
	scopeValidator ScopeValidator
}

func NewAuthorizeFlow(tokenManager TokenManager, resTypes map[string]ResponseType, scopeValidator ScopeValidator, clientStorage ClientStorage, consent Consent, session sessions.Store) *authorizeFlow {
	return &authorizeFlow{resTypes, tokenManager, clientStorage, consent, session, scopeValidator}
}

func (f *authorizeFlow) Run(ctx *AuthorizeContext) *AuthorizeRes {
	var err error
	if ctx.Challenge == "" {
		if err := f.validateAuthorize(ctx); err != nil {
			return &AuthorizeRes{Error: err}
		}

		url, err := f.createRedirectConsentURL(ctx)
		return &AuthorizeRes{RedirectURI: url, Error: err}
	}

	if err = f.validateChallenge(ctx); err != nil {
		return &AuthorizeRes{Error: err}
	}

	resType := f.resTypes[ctx.ResponseType]
	url, err := resType.GenerateUrl(ctx, f.tokenManager)
	if err != nil {
		return &AuthorizeRes{Error: err}
	}

	return &AuthorizeRes{RedirectURI: url}
}

func (f *authorizeFlow) validateAuthorize(ctx *AuthorizeContext) error {
	if ctx.ResponseType == "" {
		return errors.WithStack(errResponseTypeRequired)
	}

	if ctx.State == "" {
		return errors.WithStack(errStateRequired)
	}

	if ctx.ClientID == "" {
		return errors.WithStack(errNoClient)
	}

	if _, ok := f.resTypes[ctx.ResponseType]; !ok {
		return errors.WithStack(errResponseTypeUnSupported)
	}

	var err error
	var client *Client
	if client, err = f.clientStorage.GetClient(ctx.ClientID); err != nil {
		return err
	}
	ctx.Client = *client

	if !ctx.Client.IsValidRedirectURI(ctx.RedirectURI) {
		return errors.WithStack(errRedirectMisMatch)
	}

	return err
}

func (f *authorizeFlow) validateChallenge(ctx *AuthorizeContext) error {
	challenge, err := f.consent.ValidateChallenge(ctx.Challenge)
	if err != nil {
		return errInvalidChallenge.WithCause(err)
	}

	session, _ := f.session.Get(ctx.request, "consent")
	jdiSession, ok := session.Values["consent"].(map[string]string)
	if !ok {
		return errors.WithStack(errInvalidSession)
	}

	if jdiSession["id"] != challenge.ID {
		return errors.WithStack(errInvalidSession)
	}

	if isExpireUnix(challenge.Expired) {
		return errors.WithStack(errChallengeExpired)
	}

	ctx.ResponseType = jdiSession["resp"]
	ctx.State = jdiSession["state"]
	ctx.RedirectURI = jdiSession["redir"]
	ctx.UserID = challenge.UserID
	ctx.ClientID = challenge.ClientID
	ctx.Scopes = challenge.Scopes

	if challenge.UserID == "" {
		return errors.WithStack(errUserIDRequired)
	}

	if err := f.validateAuthorize(ctx); err != nil {
		return err
	}

	client, err := f.clientStorage.GetClient(challenge.ClientID)
	if err != nil {
		return err
	}
	ctx.Client = *client

	if f.scopeValidator != nil {
		ctx.Scopes, err = f.scopeValidator.Validate(ctx.Scopes, ctx.Client.Scopes)
		if err != nil {
			return err
		}
	}

	delete(session.Values, "consent")

	return errInternalServer.WithCause(session.Save(ctx.request, ctx.response))
}

func (f *authorizeFlow) createRedirectConsentURL(ctx *AuthorizeContext) (*url.URL, error) {
	consentUrl, consentid, err := f.consent.UrlWithChallenge(ctx.Client.ID, strings.Join(ctx.Scopes, " "))
	if err != nil {
		return nil, errInternalServer.WithCause(err)
	}

	session, _ := f.session.Get(ctx.request, "consent")

	session.Values["consent"] = map[string]string{
		"id":    consentid,
		"resp":  ctx.ResponseType,
		"redir": ctx.RedirectURI,
		"state": ctx.State,
	}

	if err = session.Save(ctx.request, ctx.response); err != nil {
		return nil, errInternalServer.WithCause(err)
	}

	return consentUrl, nil
}
