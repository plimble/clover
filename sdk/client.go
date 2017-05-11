package sdk

import "github.com/plimble/clover"

//go:generate mockery -name ClientStorage
type ClientStorage interface {
	DeleteClient(id string) error
	SaveClient(client *clover.Client) error
	GetClient(id string) (clover.Client, error)
}

type Client struct {
	storage ClientStorage
}

func NewClient(storage ClientStorage) *Client {
	return &Client{
		storage: storage,
	}
}

func (c *Client) Get(id string) (clover.Client, error) {
	return c.storage.GetClient(id)
}

func (c *Client) Delete(id string) error {
	return c.storage.DeleteClient(id)
}

func (c *Client) Save(client *clover.Client) error {
	return c.storage.SaveClient(client)
}
