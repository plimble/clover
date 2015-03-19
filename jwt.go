package clover

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	JWT_ALGO_HS256 = "HS256"
	JWT_ALGO_HS384 = "HS384"
	JWT_ALGO_HS512 = "HS512"
	JWT_ALGO_RS256 = "RS256"
	JWT_ALGO_RS384 = "RS384"
	JWT_ALGO_RS512 = "RS512"
)

func getJWTAlgorithm(algo string) jwt.SigningMethod {

	switch algo {
	case JWT_ALGO_HS256:
		return jwt.SigningMethodHS256
	case JWT_ALGO_HS384:
		return jwt.SigningMethodHS384
	case JWT_ALGO_HS512:
		return jwt.SigningMethodHS512
	case JWT_ALGO_RS256:
		return jwt.SigningMethodRS256
	case JWT_ALGO_RS384:
		return jwt.SigningMethodRS384
	case JWT_ALGO_RS512:
		return jwt.SigningMethodRS512
	}

	return jwt.SigningMethodRS256
}
