package scope_test

import (
	"testing"

	"github.com/plimble/clover/scope"
	"github.com/plimble/clover/scope/mocks"
	"github.com/stretchr/testify/require"
)

func TestScope(t *testing.T) {
	t.Run("Create", testCreate)
	t.Run("Delete", testDelete)
	t.Run("GetAll", testGetAll)
	t.Run("GetByIDs", testGetByIDs)
	t.Run("Check", testCheck)
	t.Run("HierarchicScope", testHierarchicScope)
}

var scopeSampleData = []*scope.Scope{
	{"s1", "desc"},
	{"s1:read", "desc"},
	{"s1:write", "desc"},
	{"s1:write:user", "desc"},
	{"s2", "desc"},
	{"s3", "desc"},
}

type managerTest struct {
	storage *mocks.ScopeStorage
	manager *scope.ScopeManager
}

func setup() *managerTest {
	m := &managerTest{}
	m.storage = &mocks.ScopeStorage{}
	m.manager = scope.NewManager(m.storage)

	return m
}

func testCreate(t *testing.T) {
	s := setup()

	s.storage.On("Create", scopeSampleData[0]).Return(nil)

	err := s.manager.Create("s1", "desc")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testDelete(t *testing.T) {
	s := setup()

	s.storage.On("Delete", "s1").Return(nil)

	err := s.manager.Delete("s1")
	require.NoError(t, err)
	s.storage.AssertExpectations(t)
}

func testGetAll(t *testing.T) {
	s := setup()

	s.storage.On("GetAll").Return(scopeSampleData, nil)

	scopes, err := s.manager.GetAll()
	require.NoError(t, err)
	require.Equal(t, scopeSampleData, scopes)
	s.storage.AssertExpectations(t)
}

func testGetByIDs(t *testing.T) {
	s := setup()

	s.storage.On("GetByIDs", []string{"s1", "s2"}).Return(scopeSampleData, nil)

	scopes, err := s.manager.GetByIDs([]string{"s1", "s2"})
	require.NoError(t, err)
	require.Equal(t, scopeSampleData, scopes)
	s.storage.AssertExpectations(t)
}

func testCheck(t *testing.T) {
	s := setup()
	scopes := []string{"foo", "bar.baz", "baz.baz.1", "baz.baz.2", "baz.baz.3", "baz.baz.baz"}

	usecases := []struct {
		usecase  string
		expected bool
		needles  []string
	}{
		{"Matched", true, []string{"foo.bar", "bar.baz"}},
		{"Not Matched", false, []string{"foo.bar", "bar.baz1"}},
	}

	for _, usecase := range usecases {
		t.Run(usecase.usecase, func(t *testing.T) {
			ok := s.manager.Check(scopes, usecase.needles)
			require.Equal(t, usecase.expected, ok)
		})
	}
}

func testHierarchicScope(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		var scopes = []string{}
		require.False(t, scope.HierarchicScope(scopes, "foo.bar.baz"))
		require.False(t, scope.HierarchicScope(scopes, "foo.bar"))
		require.False(t, scope.HierarchicScope(scopes, "foo"))
	})
	t.Run("case 2", func(t *testing.T) {
		scopes := []string{"foo.bar", "bar.baz", "baz.baz.1", "baz.baz.2", "baz.baz.3", "baz.baz.baz"}
		require.True(t, scope.HierarchicScope(scopes, "foo.bar.baz"))
		require.True(t, scope.HierarchicScope(scopes, "baz.baz.baz"))
		require.True(t, scope.HierarchicScope(scopes, "foo.bar"))
		require.False(t, scope.HierarchicScope(scopes, "foo"))
		require.True(t, scope.HierarchicScope(scopes, "bar.baz"))
		require.True(t, scope.HierarchicScope(scopes, "bar.baz.zad"))
		require.False(t, scope.HierarchicScope(scopes, "bar"))
		require.False(t, scope.HierarchicScope(scopes, "baz"))
	})
	t.Run("case 3", func(t *testing.T) {
		scopes := []string{"fosite.keys.create", "fosite.keys.get", "fosite.keys.delete", "fosite.keys.update"}
		require.True(t, scope.HierarchicScope(scopes, "fosite.keys.delete"))
		require.True(t, scope.HierarchicScope(scopes, "fosite.keys.get"))
		require.True(t, scope.HierarchicScope(scopes, "fosite.keys.get"))
		require.True(t, scope.HierarchicScope(scopes, "fosite.keys.update"))
	})
	t.Run("case 1", func(t *testing.T) {
		scopes := []string{"hydra", "openid", "offline"}
		require.False(t, scope.HierarchicScope(scopes, "foo.bar"))
		require.False(t, scope.HierarchicScope(scopes, "foo"))
		require.True(t, scope.HierarchicScope(scopes, "hydra"))
		require.True(t, scope.HierarchicScope(scopes, "hydra.bar"))
		require.True(t, scope.HierarchicScope(scopes, "openid"))
		require.True(t, scope.HierarchicScope(scopes, "openid.baz.bar"))
		require.True(t, scope.HierarchicScope(scopes, "offline"))
		require.True(t, scope.HierarchicScope(scopes, "offline.baz.bar.baz"))
	})
}
