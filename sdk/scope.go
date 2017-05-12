package sdk

import (
	"github.com/plimble/clover/scope"
)

//go:generate mockery -name ScopeStorage
type ScopeStorage interface {
	CreateScope(scope *scope.Scope) error
	DeleteScope(id string) error
	GetAllScope() ([]*scope.Scope, error)
}

type Scope struct {
	storage ScopeStorage
}

func NewScope(storage ScopeStorage) *Scope {
	return &Scope{storage}
}

func (s *Scope) Create(id, description string) error {
	return s.storage.CreateScope(&scope.Scope{ID: id, Description: description})
}

func (s *Scope) Delete(id string) error {
	return s.storage.DeleteScope(id)
}

func (s *Scope) GetAll() ([]*scope.Scope, error) {
	return s.storage.GetAllScope()
}
