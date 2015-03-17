package clover

func (a *AuthorizeServer) SetDefaultScopes(ids ...string) {
	a.Config.DefaultScopes = ids
}

func checkScope(available []string, request ...string) bool {
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
