package clover

type MemoryScope struct {
	scopes        map[string]*Scope
	defaultScopes []string
}

func NewMemoryScope() *MemoryScope {
	return &MemoryScope{
		scopes: make(map[string]*Scope),
	}
}

func (s *MemoryScope) SetScope(id, desc string) {
	s.scopes[id] = &Scope{id, desc}
}

func (s *MemoryScope) GetScopes(ids []string) ([]*Scope, error) {
	scopes := make([]*Scope, 0, len(ids))
	for _, id := range ids {
		if _, has := s.scopes[id]; has {
			scopes = append(scopes, s.scopes[id])
		}
	}

	return scopes, nil
}
