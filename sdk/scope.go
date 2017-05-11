package sdk

import "github.com/plimble/clover"

//go:generate mockery -name ScopeStorage
type ScopeStorage interface {
	CreateScope(scope *clover.Scope) error
	DeleteScope(id string) error
	GetAllScope() ([]*clover.Scope, error)
}

type Scope struct {
	storage ScopeStorage
}

func NewScope(storage ScopeStorage) *Scope {
	return &Scope{storage}
}

func (s *Scope) Create(id, description string) error {
	return s.storage.CreateScope(&clover.Scope{ID: id, Description: description})
}

func (s *Scope) Delete(id string) error {
	return s.storage.DeleteScope(id)
}

func (s *Scope) GetAll() ([]*clover.Scope, error) {
	return s.storage.GetAllScope()
}
