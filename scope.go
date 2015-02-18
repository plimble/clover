package clover

func (c *Clover) SetDefaultScopes(ids ...string) {
	c.Config.DefaultScopes = ids
}

func (c *Clover) GetScopeDescription(ids []string) ([]*Scope, error) {
	return c.Config.Store.GetScopes(ids)
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

func isScopeDiff(request, available []string) bool {
	vals := map[string]struct{}{}

	for _, x := range available {
		vals[x] = struct{}{}
	}

	for _, x := range request {
		if _, ok := vals[x]; !ok {
			return true
		}
	}

	return false
}
