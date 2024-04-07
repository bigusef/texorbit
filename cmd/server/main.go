package main

import (
	"context"
	"fmt"
	"github.com/bigusef/texorbit/internal/database"
	router2 "github.com/bigusef/texorbit/internal/router"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func main() {
	ctx := context.Background()
	setting := config.NewSetting()

	validate := initValidate()

	//region Database Setup
	conn, err := pgxpool.New(ctx, setting.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	queries := database.New(conn)
	//endregion

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

func initHandler(conf *config.Setting, queries *database.Queries, validate *validator.Validate) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// health check API, to make sure routers working as expected
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		util.JsonResponseWriter(w, http.StatusOK, map[string]string{"result": "OK - healthy"})
	})

	// mount all internal routers
	router.Mount("/auth", router2.NewAuthRouter())
	router.Mount("/city", router2.NewCityRouter(conf, queries, validate))

	return router
}
