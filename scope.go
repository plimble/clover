package clover

type Scope struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

//go:generate mockery -name ScopeValidator
type ScopeValidator interface {
	// Validate the request scopes and is it allowed in client scopes
	// then return the valid scopes that use can use or return error if not valid
	// example the request should available in database. then is should be subset of client scopes
	// if the request scopes is empty you may return the default scopes
	Validate(requestScopes, clientScopes []string) ([]string, error)
}
