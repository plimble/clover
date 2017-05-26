package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHierarchicScope(t *testing.T) {
	var scopes = []string{}

	assert.False(t, HierarchicScope("foo.bar.baz", scopes))
	assert.False(t, HierarchicScope("foo.bar", scopes))
	assert.False(t, HierarchicScope("foo", scopes))

	scopes = []string{"foo.bar", "bar.baz", "baz.baz.1", "baz.baz.2", "baz.baz.3", "baz.baz.baz"}
	assert.True(t, HierarchicScope("foo.bar.baz", scopes))
	assert.True(t, HierarchicScope("baz.baz.baz", scopes))
	assert.True(t, HierarchicScope("foo.bar", scopes))
	assert.False(t, HierarchicScope("foo", scopes))

	assert.True(t, HierarchicScope("bar.baz", scopes))
	assert.True(t, HierarchicScope("bar.baz.zad", scopes))
	assert.False(t, HierarchicScope("bar", scopes))
	assert.False(t, HierarchicScope("baz", scopes))

	scopes = []string{"fosite.keys.create", "fosite.keys.get", "fosite.keys.delete", "fosite.keys.update"}
	assert.True(t, HierarchicScope("fosite.keys.delete", scopes))
	assert.True(t, HierarchicScope("fosite.keys.get", scopes))
	assert.True(t, HierarchicScope("fosite.keys.get", scopes))
	assert.True(t, HierarchicScope("fosite.keys.update", scopes))

	scopes = []string{"hydra", "openid", "offline"}
	assert.False(t, HierarchicScope("foo.bar", scopes))
	assert.False(t, HierarchicScope("foo", scopes))
	assert.True(t, HierarchicScope("hydra", scopes))
	assert.True(t, HierarchicScope("hydra.bar", scopes))
	assert.True(t, HierarchicScope("openid", scopes))
	assert.True(t, HierarchicScope("openid.baz.bar", scopes))
	assert.True(t, HierarchicScope("offline", scopes))
	assert.True(t, HierarchicScope("offline.baz.bar.baz", scopes))
}
