package clover

type User struct {
	ID     string   `json:"id"`
	Scopes []string `json:"scopes"`
}

//go:generate mockery -name UserService
type UserService interface {
	GetUser(username, password string) (*User, error)
}
