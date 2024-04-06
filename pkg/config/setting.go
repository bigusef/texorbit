package config

import (
	"os"
)

type Setting struct {
	Port string
}

func NewSetting() *Setting {
	// getting server port from environment variable, and set default value if not passed from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Setting{Port: port}
}
