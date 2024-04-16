package main

import (
	"context"
	"fmt"
	"github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func main() {
	ctx := context.Background()

	// get Application Environment Settings
	setting, err := config.NewSetting()
	if err != nil {
		log.Fatal(err.Error())
	}

	validate := initValidate()

	// Database Setup
	conn := config.NewConnectionPool(ctx, setting.ConnString)
	defer conn.Close()
	queries := database.New(conn)

	handler := initHandler(setting, queries, validate)

	// start application server
	log.Println("[Restfull server] Start Store application")
	log.Println(fmt.Sprintf("Serving starting on port :%s", setting.Port))

	serverPort := fmt.Sprintf(":%s", setting.Port)
	if err := http.ListenAndServe(serverPort, handler); err != nil {
		log.Fatal(err)
	}
}

func initValidate() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
