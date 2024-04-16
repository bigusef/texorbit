package user

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type staffHandler struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

func (h *staffHandler) listStaffHandler(w http.ResponseWriter, r *http.Request) {}

func (h *staffHandler) createStaffHandler(w http.ResponseWriter, r *http.Request) {}

func (h *staffHandler) updateStaffHandler(w http.ResponseWriter, r *http.Request) {}
