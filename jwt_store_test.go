package clover

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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

var rsKey = &PublicKey{
	PublicKey:  publicKey,
	PrivateKey: privateKey,
	Algorithm:  JWT_ALGO_RS512,
}

var hsKey = &PublicKey{
	PublicKey:  hmacKey,
	PrivateKey: hmacKey,
	Algorithm:  JWT_ALGO_HS512,
}

func genJWTToken(algo jwt.SigningMethod, key string, expires int64) string {
	token := jwt.New(algo)
	token.Claims["id"] = "1111"
	token.Claims["iss"] = "2222"
	token.Claims["aud"] = "3333"
	token.Claims["sub"] = "4444"
	token.Claims["exp"] = expires
	token.Claims["iat"] = time.Now()
	token.Claims["token_type"] = "bearer"
	token.Claims["scope"] = "aaaa bbbb cccc"

	accesstoken, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	return accesstoken
}

func TestGetAccessTokenRS512(t *testing.T) {
	pubStore := NewMockPublicKeyStore()

	s := newJWTAccessTokenStore(pubStore)

	expires := time.Now().Add(time.Minute * 1).Unix()
	token := genJWTToken(jwt.SigningMethodRS512, privateKey, expires)

	expAt := &AccessToken{
		AccessToken: "1111",
		ClientID:    "3333",
		Expires:     expires,
		Scope:       []string{"aaaa", "bbbb", "cccc"},
		UserID:      "4444",
	}

	pubStore.On("GetKey", "3333").Return(rsKey, nil)

	at, err := s.GetAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, expAt, at)

}

func TestGetAccessTokenHS512(t *testing.T) {
	pubStore := NewMockPublicKeyStore()

	s := newJWTAccessTokenStore(pubStore)

	expires := time.Now().Add(time.Minute * 1).Unix()
	token := genJWTToken(jwt.SigningMethodHS512, hmacKey, expires)

	expAt := &AccessToken{
		AccessToken: "1111",
		ClientID:    "3333",
		Expires:     expires,
		Scope:       []string{"aaaa", "bbbb", "cccc"},
		UserID:      "4444",
	}

	pubStore.On("GetKey", "3333").Return(hsKey, nil)

	at, err := s.GetAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, expAt, at)

}
