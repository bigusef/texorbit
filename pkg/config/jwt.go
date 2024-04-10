package config

import (
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"os"
)

func getJWTSecret(key string) (*jwtauth.JWTAuth, error) {
	jwtSecretKey := os.Getenv(key)
	if jwtSecretKey == "" {
		return nil, errors.New(fmt.Sprintf(`messing environment variable "%s"`, key))
	}

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecretKey), nil)

	return tokenAuth, nil
}
