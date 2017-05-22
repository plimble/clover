package oauth2

type User struct {
	ID     string                 `json:"id"`
	Scopes []string               `json:"scopes"`
	Extras map[string]interface{} `json:"extras"`
}

//go:generate mockery -name UserService
type UserService interface {
	GetUser(username, password string) (*User, error)
}
