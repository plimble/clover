package scope_test

import (
	"testing"

	"github.com/plimble/clover/scope"
	"github.com/plimble/clover/scope/mocks"
	"github.com/stretchr/testify/require"
)

func TestScope(t *testing.T) {
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
	manager *scope.ScopeValidator
}

func setup() *managerTest {
	m := &managerTest{}
	m.storage = &mocks.ScopeStorage{}
	m.manager = scope.New(m.storage)

	return m
}

func testHierarchicScope(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		var scopes = []string{}
		require.False(t, scope.HierarchicScope("foo.bar.baz", scopes))
		require.False(t, scope.HierarchicScope("foo.bar", scopes))
		require.False(t, scope.HierarchicScope("foo", scopes))
	})
	t.Run("case 2", func(t *testing.T) {
		scopes := []string{"foo.bar", "bar.baz", "baz.baz.1", "baz.baz.2", "baz.baz.3", "baz.baz.baz"}
		require.True(t, scope.HierarchicScope("foo.bar.baz", scopes))
		require.True(t, scope.HierarchicScope("baz.baz.baz", scopes))
		require.True(t, scope.HierarchicScope("foo.bar", scopes))
		require.False(t, scope.HierarchicScope("foo", scopes))
		require.True(t, scope.HierarchicScope("bar.baz", scopes))
		require.True(t, scope.HierarchicScope("bar.baz.zad", scopes))
		require.False(t, scope.HierarchicScope("bar", scopes))
		require.False(t, scope.HierarchicScope("baz", scopes))
	})
	t.Run("case 3", func(t *testing.T) {
		scopes := []string{"fosite.keys.create", "fosite.keys.get", "fosite.keys.delete", "fosite.keys.update"}
		require.True(t, scope.HierarchicScope("fosite.keys.delete", scopes))
		require.True(t, scope.HierarchicScope("fosite.keys.get", scopes))
		require.True(t, scope.HierarchicScope("fosite.keys.get", scopes))
		require.True(t, scope.HierarchicScope("fosite.keys.update", scopes))
	})
	t.Run("case 1", func(t *testing.T) {
		scopes := []string{"hydra", "openid", "offline"}
		require.False(t, scope.HierarchicScope("foo.bar", scopes))
		require.False(t, scope.HierarchicScope("foo", scopes))
		require.True(t, scope.HierarchicScope("hydra", scopes))
		require.True(t, scope.HierarchicScope("hydra.bar", scopes))
		require.True(t, scope.HierarchicScope("openid", scopes))
		require.True(t, scope.HierarchicScope("openid.baz.bar", scopes))
		require.True(t, scope.HierarchicScope("offline", scopes))
		require.True(t, scope.HierarchicScope("offline.baz.bar.baz", scopes))
	})
}
