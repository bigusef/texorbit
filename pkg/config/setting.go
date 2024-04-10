package config

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"os"
)

type Setting struct {
	Port        string
	ConnString  string
	AccessAuth  *jwtauth.JWTAuth
	RefreshAuth *jwtauth.JWTAuth
}

func NewSetting() (*Setting, error) {
	var err error
	setting := &Setting{}

	// getting server port from environment variable, and set default value if not passed from env
	if setting.Port = os.Getenv("PORT"); setting.Port == "" {
		setting.Port = "8080"
	}

	// getting db url from environment variable
	if setting.ConnString = os.Getenv("DATABASE_URL"); setting.ConnString == "" {
		return nil, errors.New(`messing environment variable "DATABASE_URL"`)
	}

	// getting JWT configurations
	if setting.AccessAuth, err = getJWTSecret("JWT_ACCESS_SECRET"); err != nil {
		return nil, err
	}

	if setting.RefreshAuth, err = getJWTSecret("JWT_REFRESH_SECRET"); err != nil {
		return nil, err
	}

	return setting, nil
}
