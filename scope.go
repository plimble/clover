package clover

func (a *AuthorizeServer) SetDefaultScopes(ids ...string) {
	a.Config.DefaultScopes = ids
}

func (a *AuthorizeServer) GetScopeDescription(ids []string) ([]*Scope, error) {
	return a.Config.Store.GetScopes(ids)
}

func checkScope(request, available []string) bool {
	matched := 0

	for i := 0; i < len(request); i++ {
		for j := 0; j < len(available); j++ {
			if request[i] == available[j] {
				matched++
			}
		}
	}

	if matched != len(request) {
		return false
	}

	return true
}
