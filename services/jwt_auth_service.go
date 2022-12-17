package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"main/lib"
)

// JWTAuthService service relating to authorization
type JWTAuthService struct {
	env    lib.Env
	logger lib.Logger
}

// NewJWTAuthService creates a new auth service
func NewJWTAuthService(env lib.Env, logger lib.Logger) JWTAuthService {
	return JWTAuthService{
		env:    env,
		logger: logger,
	}
}

// Authorize authorizes the generated token
func (s JWTAuthService) Authorize(tokenString string) (bool, error) {

	keyData, err := ioutil.ReadFile("./public.key")
	if err != nil {
		return false, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return false, err
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if token.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("token expired")
		}
	}
	return false, errors.New("couldn't handle token")
}
