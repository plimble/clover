package clover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHierarchicScope(t *testing.T) {
	var scopes = []string{}
	assert.False(t, hierarchicScope(scopes, "foo.bar.baz"))
	assert.False(t, hierarchicScope(scopes, "foo.bar"))
	assert.False(t, hierarchicScope(scopes, "foo"))

	scopes = []string{"foo.bar", "bar.baz", "baz.baz.1", "baz.baz.2", "baz.baz.3", "baz.baz.baz"}
	assert.True(t, hierarchicScope(scopes, "foo.bar.baz"))
	assert.True(t, hierarchicScope(scopes, "baz.baz.baz"))
	assert.True(t, hierarchicScope(scopes, "foo.bar"))
	assert.False(t, hierarchicScope(scopes, "foo"))
	assert.True(t, hierarchicScope(scopes, "bar.baz"))
	assert.True(t, hierarchicScope(scopes, "bar.baz.zad"))
	assert.False(t, hierarchicScope(scopes, "bar"))
	assert.False(t, hierarchicScope(scopes, "baz"))

	scopes = []string{"fosite.keys.create", "fosite.keys.get", "fosite.keys.delete", "fosite.keys.update"}
	assert.True(t, hierarchicScope(scopes, "fosite.keys.delete"))
	assert.True(t, hierarchicScope(scopes, "fosite.keys.get"))
	assert.True(t, hierarchicScope(scopes, "fosite.keys.get"))
	assert.True(t, hierarchicScope(scopes, "fosite.keys.update"))

	scopes = []string{"hydra", "openid", "offline"}
	assert.False(t, hierarchicScope(scopes, "foo.bar"))
	assert.False(t, hierarchicScope(scopes, "foo"))
	assert.True(t, hierarchicScope(scopes, "hydra"))
	assert.True(t, hierarchicScope(scopes, "hydra.bar"))
	assert.True(t, hierarchicScope(scopes, "openid"))
	assert.True(t, hierarchicScope(scopes, "openid.baz.bar"))
	assert.True(t, hierarchicScope(scopes, "offline"))
	assert.True(t, hierarchicScope(scopes, "offline.baz.bar.baz"))
}
