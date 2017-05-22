package oauth2

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

func (c *Client) IsValidRedirectURI(uri string) bool {
	for _, u := range c.RedirectURIs {
		if u == uri {
			return true
		}
	}

	return false
}
