package memory

import (
	"sync"

	"github.com/plimble/clover"
)

type userService struct {
	sync.Mutex
	user map[string]*clover.User
}

func NewUserService() clover.UserService {
	return &userService{
		user: make(map[string]*clover.User),
	}
}

func (s *userService) Flush() {
	s.Lock()
	s.user = make(map[string]*clover.User)
	s.Unlock()
}

func (s *userService) GetUser(username, password string) (*clover.User, error) {
	s.Lock()
	defer s.Unlock()
	user, ok := s.user[username+":"+password]
	if !ok {
		return nil, errNotFound
	}

	return user, nil
}
