package user

import (
	"encoding/json"
	"errors"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/middleware"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type staffHandler struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

func (h *staffHandler) listStaffHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := ctx.Value("pagination").(*middleware.Paginator)

	staff, err := h.queries.AllStaff(ctx, db.AllStaffParams{page.Limit, page.Offset})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result := make([]*listStaff, len(staff))
	for i, v := range staff {
		result[i] = &listStaff{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			JoinDate: v.JoinDate.Time,
		}
	}

	count, err := h.queries.AllStaffCount(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, result, count)
}

func (h *staffHandler) createStaffHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input newStaff
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(input); err != nil {
		ts := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			ts[err.Field()] = err.Tag()
		}

		util.JsonResponseWriter(w, http.StatusBadRequest, ts)
		return
	}

	if _, err := h.queries.GetUserByEmail(ctx, input.Email); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		util.JsonResponseWriter(w, http.StatusBadRequest, map[string]string{"email": "email already used by another user"})
		return
	}

	user, err := h.queries.CreateUser(ctx, db.CreateUserParams{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: pgtype.Text{String: input.PhoneNumber, Valid: true},
		IsStaff:     true,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	util.JsonResponseWriter(w, http.StatusCreated, staffInfo{
		Id:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Avatar:      user.Avatar.String,
		PhoneNumber: user.PhoneNumber.String,
		Status:      string(user.Status),
		JoinDate:    user.JoinDate.Time,
		LastLogin:   user.LastLogin.Time,
	})
}

func (h *staffHandler) updateStaffHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUSerById(ctx, id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var input updateStaff
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(input); err != nil {
		ts := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			ts[err.Field()] = err.Tag()
		}

		util.JsonResponseWriter(w, http.StatusBadRequest, ts)
		return
	}

	if input.Email != user.Email {
		if _, err := h.queries.GetUserByEmail(ctx, input.Email); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			util.JsonResponseWriter(w, http.StatusBadRequest, map[string]string{"email": "email already used by another user"})
			return
		}
	}

	updatedUser, err := h.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:          id,
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: pgtype.Text{String: input.PhoneNumber, Valid: true},
		Status:      input.Status,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	util.JsonResponseWriter(w, http.StatusOK, staffInfo{
		Id:          updatedUser.ID,
		Name:        updatedUser.Name,
		Email:       updatedUser.Email,
		Avatar:      updatedUser.Avatar.String,
		PhoneNumber: updatedUser.PhoneNumber.String,
		Status:      string(updatedUser.Status),
		JoinDate:    updatedUser.JoinDate.Time,
		LastLogin:   updatedUser.LastLogin.Time,
	})
}
