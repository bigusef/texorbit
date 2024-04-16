package user

import (
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/go-playground/validator/v10"
)

type customerHandler struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

/* customers
- get profile data
- update profile data
*/
