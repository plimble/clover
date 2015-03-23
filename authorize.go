package clover

import (
	"strings"
)

type authController struct {
	config        *AuthConfig
	authRespType  AuthRespType
	tokenRespType AuthRespType
}

type authorizeRequest struct {
	responseType string
	scope        string
	clientID     string
	redirectURI  string
	state        string
	client       Client
	scopeArr     []string
	respType     AuthRespType
}

func newAuthController(config *AuthConfig, authRespType, tokenRespType AuthRespType) *authController {
	return &authController{config, authRespType, tokenRespType}
}

func (a *authController) handleAuthorize(ar *authorizeRequest, isAuthorized bool) *Response {
	// the user declined access to the client's application
	if !isAuthorized {
		return errUserDeniedAccess.setRedirect(ar.redirectURI, ar.responseType, ar.state)
	}

	//validate request
	resp := a.validateAuthRequest(ar)
	if resp != nil {
		return resp
	}

	return ar.respType.GetAuthResponse(ar, ar.client, ar.scopeArr)
}

func (a *authController) validateAuthRequest(ar *authorizeRequest) *Response {
	var resp *Response

	if resp = a.validateAuthClientID(ar); resp != nil {
		return resp
	}

	if resp = a.validateAuthRedirectURI(ar); resp != nil {
		return resp
	}

	if resp = a.validateAuthState(ar); resp != nil {
		return resp
	}

	if resp = a.validateAuthRespType(ar); resp != nil {
		resp.setRedirect(ar.redirectURI, ar.responseType, ar.state)
		return resp
	}

	if resp = a.validateAuthScope(ar); resp != nil {
		resp.setRedirect(ar.redirectURI, ar.responseType, ar.state)
		return resp
	}

	return nil
}

func (a *authController) validateAuthClientID(ar *authorizeRequest) *Response {
	var err error

	if ar.clientID == "" {
		return errNoClientID
	}

	//get client
	ar.client, err = a.config.AuthServerStore.GetClient(ar.clientID)
	if err != nil {
		return errInvalidClientID
	}

	return nil
}

func (a *authController) validateAuthState(ar *authorizeRequest) *Response {
	if a.config.StateParamRequired && ar.state == "" {
		return errStateRequired
	}

	return nil
}

// Make sure a valid redirect_uri was supplied. If specified, it must match the clientData URI.
func (a *authController) validateAuthRedirectURI(ar *authorizeRequest) *Response {
	if ar.redirectURI == "" {
		if ar.client.GetRedirectURI() == "" {
			return errNoRedirectURI
		}
		ar.redirectURI = ar.client.GetRedirectURI()
		return nil
	}

	if ar.redirectURI != "" && ar.client.GetRedirectURI() != "" && ar.client.GetRedirectURI() != ar.redirectURI {
		return errRedirectMismatch
	}

	return nil
}

func (a *authController) validateAuthRespType(ar *authorizeRequest) *Response {
	switch ar.responseType {
	case "":
		return errInvalidRespType
	case RESP_TYPE_CODE:
		if !checkGrantType(ar.client.GetGrantType(), AUTHORIZATION_CODE) {
			return errUnAuthorizedGrant
		}
		ar.respType = a.authRespType
	case RESP_TYPE_TOKEN:
		if !a.config.AllowImplicit {
			return errUnSupportedImplicit
		}

		if !checkGrantType(ar.client.GetGrantType(), IMPLICIT) {
			return errUnAuthorizedGrant
		}
		ar.respType = a.tokenRespType
	default:
		return errCodeUnSupportedGrant
	}

	return nil
}

func (a *authController) validateAuthScope(ar *authorizeRequest) *Response {
	scopes := strings.Fields(ar.scope)

	if len(scopes) == 0 {
		ar.scopeArr = a.config.DefaultScopes
		return nil
	}

	if len(ar.client.GetScope()) == 0 || !checkScope(ar.client.GetScope(), scopes...) {
		return errUnSupportedScope
	}

	ar.scopeArr = scopes
	return nil
}
