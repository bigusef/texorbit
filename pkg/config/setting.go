package config

import (
	"log"
	"os"
)

type Setting struct {
	Port       string
	ConnString string
}

func NewSetting() *Setting {
	// getting server port from environment variable, and set default value if not passed from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// getting db url from environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal(`messing environment variable "DATABASE_URL", terminate server....`)
	}

	return &Setting{Port: port, ConnString: dsn}
}
