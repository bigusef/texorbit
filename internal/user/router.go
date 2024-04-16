package user

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func AuthRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	h := &authHandler{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	// public
	r.Post("/login", h.login)
	r.Post("/staff-login", h.staffLogin)

	// staff and customers
	r.With(
		jwtauth.Verifier(conf.RefreshAuth),
		jwtauth.Authenticator(conf.RefreshAuth),
	).Get("/refresh", h.refreshAccessToken)

	return r
}

func StaffRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	h := &staffHandler{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	return r
}

func CustomerRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	h := &customerHandler{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	return r
}
