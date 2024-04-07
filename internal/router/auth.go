package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewAuthRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/register", registrationHandler)
	router.Post("/login", loginHandler)
	router.Post("/login/admin", adminLoginHandler)
	router.Post("/forget-password", forgetPasswordHandler)
	router.Post("/reset-password", resetPasswordHandler)
	router.Get("/logout", logoutHandler)
	router.Get("/me", getUserDataHandler)
	router.Put("/me", updateUserDataHandler)

	return router
}

func registrationHandler(w http.ResponseWriter, r *http.Request)   {}
func loginHandler(w http.ResponseWriter, r *http.Request)          {}
func adminLoginHandler(w http.ResponseWriter, r *http.Request)     {}
func logoutHandler(w http.ResponseWriter, r *http.Request)         {}
func forgetPasswordHandler(w http.ResponseWriter, r *http.Request) {}
func resetPasswordHandler(w http.ResponseWriter, r *http.Request)  {}
func getUserDataHandler(w http.ResponseWriter, r *http.Request)    {}
func updateUserDataHandler(w http.ResponseWriter, r *http.Request) {}
