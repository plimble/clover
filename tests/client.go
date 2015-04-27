package tests

import (
	"github.com/plimble/rand"
)

type TestClient struct {
	ClientID     string
	ClientSecret string
	GrantType    []string
	UserID       string
	Scope        []string
	RedirectURI  string
	Data         map[string]interface{}
}

func (c *TestClient) GetClientID() string {
	return c.ClientID
}

func (c *TestClient) GetClientSecret() string {
	return c.ClientSecret
}

func (c *TestClient) GetGrantType() []string {
	return c.GrantType
}

func (c *TestClient) GetUserID() string {
	return c.UserID
}

func (c *TestClient) GetScope() []string {
	return c.Scope
}

func (c *TestClient) GetRedirectURI() string {
	return c.RedirectURI
}

func (c *TestClient) GetData() map[string]interface{} {
	return c.Data
}

func genTestClient() *TestClient {
	return &TestClient{
		ClientID:     rand.Digits(5),
		ClientSecret: rand.Digits(5),
		UserID:       rand.Digits(3),
		RedirectURI:  "http://www.localhost",
	}
}
