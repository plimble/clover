package sdk

import (
	"testing"

	"github.com/plimble/clover"
	"github.com/plimble/clover/sdk/mocks"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("Get", testClientGet)
	t.Run("Delete", testClientDelete)
	t.Run("Save", testClientSave)
}

var clientSampleData = []clover.Client{
	{"c1", "client", "secret", []string{"http://example.com"}, []string{"password"}, []string{"s1", "s2"}, false, "time", false},
}

type clientTest struct {
	store  *mocks.ClientStorage
	client *Client
}

func setup() *clientTest {
	c := &clientTest{}
	c.store = &mocks.ClientStorage{}
	c.client = NewClient(c.store)

	return c
}

func testClientGet(t *testing.T) {
	s := setup()

	s.store.On("GetClient", "c1").Return(clientSampleData[0], nil)

	client, err := s.client.Get("c1")
	require.NoError(t, err)
	require.Equal(t, clientSampleData[0], client)
	s.store.AssertExpectations(t)
}

func testClientDelete(t *testing.T) {
	s := setup()

	s.store.On("DeleteClient", "c1").Return(nil)

	err := s.client.Delete("c1")
	require.NoError(t, err)
	s.store.AssertExpectations(t)
}

func testClientSave(t *testing.T) {
	s := setup()

	s.store.On("SaveClient", &clientSampleData[0]).Return(nil)

	err := s.client.Save(&clientSampleData[0])
	require.NoError(t, err)
	s.store.AssertExpectations(t)
}
