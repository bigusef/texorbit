package main

import (
	"github.com/bigusef/texorbit/internal/city"
	"github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/internal/user"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func initHandler(conf *config.Setting, queries *database.Queries, validate *validator.Validate) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// health check API, to make sure routers working as expected
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		util.JsonResponseWriter(w, http.StatusOK, map[string]string{"result": "OK - healthy"})
	})

	// mount all internal routers
	router.Mount("/auth", user.AuthRouter(conf, queries, validate))
	router.Mount("/staff", user.StaffRouter(conf, queries, validate))
	router.Mount("/customer", user.CustomerRouter(conf, queries, validate))
	router.Mount("/city", city.NewRouter(conf, queries, validate))

	return router
}
