package router

import (
	"github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func New(conf *config.Setting, queries *database.Queries) http.Handler {
	router := chi.NewRouter()
	setupMiddleware(router)

	// health check API, to make sure routers working as expected
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		util.JsonResponseWriter(w, http.StatusOK, map[string]string{"result": "OK - healthy"})
	})

	// mount all internal routers
	router.Mount("/city", newCityRouter(conf, queries))

	return router
}

func setupMiddleware(mux *chi.Mux) {
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}
