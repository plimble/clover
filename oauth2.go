package clover

import (
	"net/http"

	"github.com/gorilla/sessions"
)

//go:generate mockery -name OAuth2
type OAuth2 interface {
	IntrospectionHandler(w http.ResponseWriter, r *http.Request)
	AccessTokenHandler(w http.ResponseWriter, r *http.Request)
	AuthorizeHandler(w http.ResponseWriter, r *http.Request)
	RevokeHandler(w http.ResponseWriter, r *http.Request)
}

type oauth2 struct {
	clientStorage     ClientStorage
	tokenStorage      TokenStorage
	strategy          *Strategy
	tokenManager      TokenManager
	consent           Consent
	authorizeFlow     *authorizeFlow
	accessTokenFlow   *accessTokenFlow
	introspectionFlow *introspectionFlow
	revokeFlow        *revokeFlow
}

func New(clientStorage ClientStorage, tokenStorage TokenStorage, strategy *Strategy) OAuth2 {
	o := &oauth2{
		clientStorage: clientStorage,
		tokenStorage:  tokenStorage,
		strategy:      strategy,
	}

	o.tokenManager = &tokenManager{
		accessTokenGenerator:  o.strategy.accessTokenGenerator,
		refreshTokenGenerator: o.strategy.refreshTokenGenerator,
		authCodeGenerator:     o.strategy.authorizeTokenGenerator,
		tokenStore:            o.tokenStorage,
	}

	o.consent = NewConsent(
		nil,
		o.strategy.authorizeConfig.ConsentUrl,
		o.strategy.authorizeConfig.ChallengeLifeSpan,
	)

	o.accessTokenFlow = &accessTokenFlow{
		tokenManager:   o.tokenManager,
		grantTypes:     o.strategy.grantTypes,
		scopeValidator: o.strategy.scopeValidator,
		clientStorage:  o.clientStorage,
	}

	o.authorizeFlow = &authorizeFlow{
		resTypes:       o.strategy.authorizeConfig.resTypes,
		tokenManager:   o.tokenManager,
		clientStorage:  o.clientStorage,
		consent:        o.consent,
		session:        sessions.NewCookieStore([]byte(o.strategy.authorizeConfig.CookieSecret)),
		scopeValidator: o.strategy.scopeValidator,
	}

	o.introspectionFlow = &introspectionFlow{
		clientStorage: o.clientStorage,
		tokenManager:  o.tokenManager,
	}

	o.revokeFlow = &revokeFlow{
		clientStorage: o.clientStorage,
		tokenManager:  o.tokenManager,
	}

	return o
}

func (o *oauth2) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseAccessTokenRequest(w, r)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	res, err := o.accessTokenFlow.Run(ctx)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	writeJson(&ctx.HTTPContext, 200, res)
}

func (o *oauth2) IntrospectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseIntospectionRequest(w, r)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	res, err := o.introspectionFlow.Run(ctx)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	writeJson(&ctx.HTTPContext, 200, res)
}

func (o *oauth2) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseAuthorizeRequest(w, r)
	if err != nil {
		redirectError(&ctx.HTTPContext, o.strategy.authorizeConfig.ConsentUrl, err)
		return
	}

	res := o.authorizeFlow.Run(ctx)
	if res.Error != nil {
		redirectError(&ctx.HTTPContext, o.strategy.authorizeConfig.ConsentUrl, err)
		return
	}

	redirect(&ctx.HTTPContext, res.RedirectURI.String())
}

func (o *oauth2) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseRevokeRequest(w, r)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	err = o.revokeFlow.Run(ctx)
	if err != nil {
		writeJsonError(&ctx.HTTPContext, err)
		return
	}

	writeJson(&ctx.HTTPContext, 200, nil)
}
