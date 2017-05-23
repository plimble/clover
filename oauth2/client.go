package oauth2

type Client struct {
	ID           string   `json:"id"`
	Name         string   `json:"n"`
	Secret       string   `json:"s"`
	RedirectURIs []string `json:"rdr"`
	GrantTypes   []string `json:"gt"`
	Scopes       []string `json:"scp"`
	Public       bool     `json:"pub"`
	CreatedAt    string   `json:"cat"`
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
