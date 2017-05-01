package clover

type Client struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Secret       string   `json:"secret"`
	RedirectURIs []string `json:"redirect_uris"`
	GrantTypes   []string `json:"grant_types"`
	Scopes       []string `json:"scopes"`
	Public       bool     `json:"public"`
	CreatedAt    string   `json:"created_at"`
}

func (c *Client) IsGrantType(grant string) bool {
	for _, g := range c.GrantTypes {
		if g == grant {
			return true
		}
	}

	return false
}

func (c *Client) IsRedirectURI(uri string) bool {
	for _, u := range c.RedirectURIs {
		if u == uri {
			return true
		}
	}

	return false
}

//go:generate mockery -name ClientManager -case underscore
type ClientManager interface {
	GetClient(id, secret string) (*Client, error)
}

type clientManager struct {
	storage ClientStore
}

func NewClientManager(storage ClientStore) ClientManager {
	return &clientManager{
		storage: storage,
	}
}

func (c *clientManager) GetClient(id, secret string) (*Client, error) {
	return nil, nil
}
