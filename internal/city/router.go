package city

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func NewRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	h := &cityHandler{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	//only staff
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.AccessAuth))
		r.Use(jwtauth.Authenticator(conf.AccessAuth))
		r.Use(middleware.StaffPermission)

		r.Post("/", h.createCity)
		r.With(middleware.Pagination).Get("/", h.listCities)
		r.Put("/{id}", h.updateCity)
		r.Delete("/{id}", h.deleteCity)
	})

	//public
	r.With(middleware.Pagination).Get("/active", h.listActiveCities)

	return r
}
