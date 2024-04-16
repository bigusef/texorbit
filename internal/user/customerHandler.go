package user

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type customerHandler struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

func (h *customerHandler) getUserInfo(w http.ResponseWriter, r *http.Request) {}

func (h *customerHandler) updateUserInfo(w http.ResponseWriter, r *http.Request) {}

func (h *customerHandler) listAllCustomers(w http.ResponseWriter, r *http.Request) {}

func (h *customerHandler) getCustomerInfo(w http.ResponseWriter, r *http.Request) {}

func (h *customerHandler) updateCustomerInfo(w http.ResponseWriter, r *http.Request) {}
