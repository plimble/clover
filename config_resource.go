package clover

type ResourceConfig struct {
	WWWRealm        string
	AuthServerStore AuthServerStore
	PublicKeyStore  PublicKeyStore
}

func NewResourceConfig(authStore AuthServerStore) *ResourceConfig {
	return &ResourceConfig{
		WWWRealm:        "Service",
		AuthServerStore: authStore,
	}
}

func (c *ResourceConfig) UseJWTAccessTokens(store PublicKeyStore) {
	c.PublicKeyStore = store
	c.AuthServerStore = newJWTTokenStore(store)
}
