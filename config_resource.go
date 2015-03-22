package clover

type ResourceConfig struct {
	AuthServerStore AuthServerStore
	PublicKeyStore  PublicKeyStore
}

func NewResourceConfig(authStore AuthServerStore) *ResourceConfig {
	return &ResourceConfig{}
}

func (c *ResourceConfig) UseJWTAccessTokens(store PublicKeyStore) {
	c.PublicKeyStore = store
	c.AuthServerStore = newJWTAccessTokenStore(store)
}
