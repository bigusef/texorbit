package router

import (
	"encoding/json"
	"errors"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

type userInputData struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Avatar string `json:"avatar" validate:"required,url"`
}

type userRouter struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

func NewAuthRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	router := chi.NewRouter()
	handler := &userRouter{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	// public
	router.Post("/login", handler.login)
	router.Post("/staff-login", handler.staffLogin)

	// staff and customers
	router.With(
		jwtauth.Verifier(conf.RefreshAuth),
		jwtauth.Authenticator(conf.RefreshAuth),
	).Get("/refresh", handler.refreshAccessToken)

	return router
}

func (h *userRouter) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//TODO: change here to get the data from oauth2 logic
	var payload userInputData
	createAccount := false

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		ts := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			ts[err.Field()] = err.Tag()
		}

		util.JsonResponseWriter(w, http.StatusBadRequest, ts)
		return
	}

	// start of logic get or create user
	user, err := h.queries.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			util.ErrorResponseWriter(w, http.StatusInternalServerError, "issue in getting user data")
			return
		}

		// when the user ont exist in our DB then change create user flag to true
		createAccount = true
	}

	if createAccount {
		newUser, err := h.queries.CreateUser(ctx, db.CreateUserParams{
			Name:    payload.Name,
			Email:   payload.Email,
			Avatar:  pgtype.Text{String: payload.Avatar, Valid: true},
			IsStaff: false,
		})
		if err != nil {
			util.ErrorResponseWriter(w, http.StatusBadRequest, "failed to create user")
			return
		}

		user = newUser
	}

	// validate user not blocked
	if !user.IsActive() {
		util.ErrorResponseWriter(w, http.StatusForbidden, "There are issue in your account, please contact with support.")
	}

	// update user data
	// TODO: complete this
	//  - update last login
	//  - update user data from payload

	// get access token and refresh token
	_, accessToken, _ := h.conf.AccessAuth.Encode(
		map[string]interface{}{
			"sub":   user.ID.String(),
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
			"staff": user.IsStaff,
		},
	)
	_, refreshToken, _ := h.conf.RefreshAuth.Encode(
		map[string]interface{}{
			"sub": user.ID.String(),
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		},
	)

	response := struct {
		Name         string
		Email        string      `json:"email"`
		PhoneNumber  pgtype.Text `json:"phone_number"`
		Avatar       pgtype.Text `json:"avatar"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}{
		user.Name,
		user.Email,
		user.PhoneNumber,
		user.Avatar,
		accessToken,
		refreshToken,
	}

	util.JsonResponseWriter(w, http.StatusOK, response)
}

func (h *userRouter) staffLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//TODO: change here to get the data from oauth2 logic
	var payload userInputData

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		ts := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			ts[err.Field()] = err.Tag()
		}

		util.JsonResponseWriter(w, http.StatusBadRequest, ts)
		return
	}

	// start of logic get or create user
	user, err := h.queries.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.ErrorResponseWriter(w, http.StatusNotFound, "this user does not exist in the system.")
			return
		}

		util.ErrorResponseWriter(w, http.StatusInternalServerError, "issue in getting user data")
		return
	}

	if !user.IsActive() || !user.IsStaff {
		util.ErrorResponseWriter(w, http.StatusForbidden, "There are issue in your account please contact with your IT support.")
		return
	}

	// get access token and refresh token
	_, accessToken, accessErr := h.conf.AccessAuth.Encode(
		map[string]interface{}{
			"sub":   user.ID.String(),
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
			"staff": user.IsStaff,
		},
	)
	_, refreshToken, refreshErr := h.conf.RefreshAuth.Encode(
		map[string]interface{}{
			"sub": user.ID.String(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	if accessErr != nil || refreshErr != nil {
		util.ErrorResponseWriter(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response := struct {
		Name         string
		Email        string      `json:"email"`
		PhoneNumber  pgtype.Text `json:"phone_number"`
		Avatar       pgtype.Text `json:"avatar"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}{
		user.Name,
		user.Email,
		user.PhoneNumber,
		user.Avatar,
		accessToken,
		refreshToken,
	}

	util.JsonResponseWriter(w, http.StatusOK, response)
}

func (h *userRouter) refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get userId from refresh token
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub := claims["sub"].(string)
	userId, err := uuid.Parse(sub)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// get user from DB
	user, err := h.queries.GetUSerById(ctx, userId)
	if err != nil {
		util.JsonResponseWriter(w, http.StatusNotFound, "this user does not exist in the system.")
		return
	}

	// user validation
	if !user.IsActive() {
		util.ErrorResponseWriter(w, http.StatusForbidden, "There are issue in your account, please contact with support.")
		return
	}

	// generate new access token
	_, accessToken, _ := h.conf.AccessAuth.Encode(
		map[string]interface{}{
			"sub":   user.ID.String(),
			"exp":   time.Now().Add(time.Minute * 15).Unix(),
			"staff": user.IsStaff,
		},
	)

	util.JsonResponseWriter(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
