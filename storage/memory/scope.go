package memory

import (
	"sync"

	"github.com/plimble/clover/scope"
)

type ScopeStorage struct {
	sync.Mutex
	scope map[string]*scope.Scope
}

func NewScopeStorage() *ScopeStorage {
	return &ScopeStorage{
		scope: make(map[string]*scope.Scope),
	}
}

func (s *ScopeStorage) Flush() {
	s.Lock()
	s.scope = make(map[string]*scope.Scope)
	s.Unlock()
}

func (s *ScopeStorage) CreateScope(scope *scope.Scope) error {
	s.Lock()
	defer s.Unlock()
	s.scope[scope.ID] = scope
	return nil
}

func (s *ScopeStorage) GetScopeByIDs(ids []string) ([]*scope.Scope, error) {
	s.Lock()
	defer s.Unlock()
	scopes := []*scope.Scope{}

	for _, id := range ids {
		scope, ok := s.scope[id]
		if ok {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

func (s *ScopeStorage) DeleteScope(id string) error {
	s.Lock()
	defer s.Unlock()
	_, ok := s.scope[id]
	if !ok {
		return errNotFound
	}

	delete(s.scope, id)

	return nil
}

func (s *ScopeStorage) GetAllScope() ([]*scope.Scope, error) {
	s.Lock()
	defer s.Unlock()

	scopes := make([]*scope.Scope, len(s.scope))
	index := 0
	for _, scope := range s.scope {
		scopes[index] = scope
		index++
	}

	return scopes, nil
}
