package memory

import (
	"errors"
	"github.com/plimble/clover"
)

var (
	errNotFound = errors.New("not found")
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserFunc func(username, password string) (string, []string, error)
type GetClientFunc func(clientID string) (clover.Client, error)

type Storage struct {
	Client      map[string]*clover.DefaultClient
	Refresh     map[string]*clover.RefreshToken
	AuthCode    map[string]*clover.AuthorizeCode
	AccessToken map[string]*clover.AccessToken
	User        map[string]*User
	PublicKey   map[string]*clover.PublicKey
}

func New() *Storage {
	return &Storage{
		Client:      make(map[string]*clover.DefaultClient),
		Refresh:     make(map[string]*clover.RefreshToken),
		AuthCode:    make(map[string]*clover.AuthorizeCode),
		AccessToken: make(map[string]*clover.AccessToken),
		User:        make(map[string]*User),
		PublicKey:   make(map[string]*clover.PublicKey),
	}
}

func (s *Storage) GetClient(id string) (clover.Client, error) {
	client, ok := s.Client[id]
	if !ok {
		return nil, errNotFound
	}

	return client, nil
}

func (s *Storage) SetAccessToken(accessToken *clover.AccessToken) error {
	s.AccessToken[accessToken.AccessToken] = accessToken
	return nil
}

func (s *Storage) GetAccessToken(at string) (*clover.AccessToken, error) {
	accesstoken, ok := s.AccessToken[at]
	if !ok {
		return nil, errNotFound
	}

	return accesstoken, nil
}

func (s *Storage) SetRefreshToken(rt *clover.RefreshToken) error {
	s.Refresh[rt.RefreshToken] = rt
	return nil
}

func (s *Storage) GetRefreshToken(rt string) (*clover.RefreshToken, error) {
	refreshtoken, ok := s.Refresh[rt]
	if !ok {
		return nil, errNotFound
	}

	return refreshtoken, nil
}

func (s *Storage) RemoveRefreshToken(rt string) error {
	_, ok := s.Refresh[rt]
	if !ok {
		return errNotFound
	}

	delete(s.Refresh, rt)

	return nil
}

func (s *Storage) SetAuthorizeCode(ac *clover.AuthorizeCode) error {
	s.AuthCode[ac.Code] = ac
	return nil
}

func (s *Storage) GetAuthorizeCode(code string) (*clover.AuthorizeCode, error) {
	authcode, ok := s.AuthCode[code]
	if !ok {
		return nil, errNotFound
	}

	return authcode, nil
}

func (s *Storage) GetUser(username, password string) (string, []string, error) {
	user, ok := s.User[username]
	if !ok {
		return "", nil, errNotFound
	}

	if password != user.Password {
		return "", nil, errors.New("invalid username or password")
	}

	return user.Username, nil, nil
}

func (s *Storage) GetKey(clientID string) (*clover.PublicKey, error) {
	if clientID == "" {
		return &clover.PublicKey{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
			Algorithm:  clover.JWT_ALGO_RS512,
		}, nil
	}

	pub, ok := s.PublicKey[clientID]
	if !ok {
		return nil, errNotFound
	}

	return pub, nil
}

func (s *Storage) AddRSKey(clientID string) {
	s.PublicKey[clientID] = &clover.PublicKey{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Algorithm:  clover.JWT_ALGO_RS512,
	}
}

func (s *Storage) AddHSKey(clientID string) {
	s.PublicKey[clientID] = &clover.PublicKey{
		PublicKey:  hmacKey,
		PrivateKey: hmacKey,
		Algorithm:  clover.JWT_ALGO_HS512,
	}
}

var privateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----
`

var publicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`

var hmacKey = `#5K+･ｼミew{ｦ住ｳ(跼Tﾉ(ｩ┫ﾒP.ｿﾓ燾辻G�感ﾃwb="=.!r.Oﾀﾍ奎gﾐ｣`
