package user

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/middleware"
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
	r.Use(jwtauth.Verifier(conf.AccessAuth))
	r.Use(jwtauth.Authenticator(conf.AccessAuth))
	r.Use(middleware.StaffPermission)

	// Only Staff users [admin]
	r.Get("/", h.listStaffHandler)
	r.Post("/", h.createStaffHandler)
	r.Put("/{id}", h.updateStaffHandler)

	return r
}

func CustomerRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	h := &customerHandler{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	r.Use(jwtauth.Verifier(conf.AccessAuth))
	r.Use(jwtauth.Authenticator(conf.AccessAuth))

	// only authenticated user will get this based on auth token
	r.Get("/me", h.getUserInfo)
	r.Put("/me", h.updateUserInfo)

	// only staff users [Admin]
	r.Group(func(ir chi.Router) {
		ir.Use(middleware.StaffPermission)

		ir.Get("/", h.listAllCustomers)
		ir.Get("/{id}", h.getCustomerInfo)
		ir.Put("/{id}", h.updateCustomerInfo)
	})

	return r
}
