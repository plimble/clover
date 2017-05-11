package memory

import (
	"errors"
	"sync"

	"github.com/plimble/clover"
)

var (
	errNotFound = errors.New("not found")
)

type clientStorage struct {
	sync.Mutex
	client map[string]*clover.Client
}

func NewClientStorage() clover.ClientStorage {
	return &clientStorage{
		client: make(map[string]*clover.Client),
	}
}

func (s *clientStorage) Flush() {
	s.Lock()
	s.client = make(map[string]*clover.Client)
	s.Unlock()
}

func (s *clientStorage) DeleteClient(id string) error {
	s.Lock()
	defer s.Unlock()
	_, ok := s.client[id]
	if !ok {
		return errNotFound
	}

	delete(s.client, id)

	return nil
}

func (s *clientStorage) SaveClient(client *clover.Client) error {
	s.Lock()
	defer s.Unlock()
	s.client[client.ID] = client
	return nil
}

func (s *clientStorage) GetClientWithSecret(id, secret string) (*clover.Client, error) {
	s.Lock()
	defer s.Unlock()
	client, ok := s.client[id]
	if !ok {
		return nil, errNotFound
	}

	if client.Secret != secret {
		return nil, errors.New("secret mismatch")
	}

	return client, nil
}

func (s *clientStorage) GetClient(id string) (*clover.Client, error) {
	s.Lock()
	defer s.Unlock()
	client, ok := s.client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}
