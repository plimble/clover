package tests

type TestUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Data     map[string]interface{}
	Scope    []string
}

func (u *TestUser) GetID() string {
	return u.ID
}

func (u *TestUser) GetUsername() string {
	return u.Username
}

func (u *TestUser) GetPassword() string {
	return u.Password
}

func (u *TestUser) GetData() map[string]interface{} {
	return u.Data
}

func (u *TestUser) GetScope() []string {
	return u.Scope
}
