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

/* staff's users
- list all staff with filter
- create staff
- update staff
- list customers with filter
- update customers
*/

func (h *staffHandler) listStaffHandler(w http.ResponseWriter, r *http.Request) {}
