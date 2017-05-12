package sdk

import (
	"testing"

	"github.com/plimble/clover/scope"
	"github.com/plimble/clover/sdk/mocks"
	"github.com/stretchr/testify/require"
)

func TestScope(t *testing.T) {
	t.Run("Create", testScopeCreate)
	t.Run("Delete", testScopeDelete)
	t.Run("GetAll", testScopeGetAll)
}

var scopeSampleData = []*scope.Scope{
	{"s1", "desc"},
	{"s1:read", "desc"},
	{"s1:write", "desc"},
	{"s1:write:user", "desc"},
	{"s2", "desc"},
	{"s3", "desc"},
}

type scopeTest struct {
	storage *mocks.ScopeStorage
	scope   *Scope
}

func setupScope() *scopeTest {
	s := &scopeTest{}
	s.storage = &mocks.ScopeStorage{}
	s.scope = NewScope(s.storage)

	return s
}

func testScopeCreate(t *testing.T) {
	s := setupScope()

	s.storage.On("CreateScope", scopeSampleData[0]).Return(nil)

	err := s.scope.Create("s1", "desc")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testScopeDelete(t *testing.T) {
	s := setupScope()

	s.storage.On("DeleteScope", "s1").Return(nil)

	err := s.scope.Delete("s1")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testScopeGetAll(t *testing.T) {
	s := setupScope()

	s.storage.On("GetAllScope").Return(scopeSampleData, nil)

	scopes, err := s.scope.GetAll()
	require.NoError(t, err)
	require.Equal(t, scopeSampleData, scopes)
	s.storage.AssertExpectations(t)
}
