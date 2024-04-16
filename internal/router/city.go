package router

import (
	"encoding/json"
	"errors"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/middleware"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

// Schema data
type cityRouter struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

type cityInput struct {
	NameEn   string `json:"name_en" validate:"required"`
	NameAr   string `json:"name_ar" validate:"required"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type cityResponse struct {
	ID       int64  `json:"id"`
	NameEn   string `json:"name_en"`
	NameAr   string `json:"name_ar"`
	IsActive bool   `json:"is_active"`
}

type activeCityResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewCityRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	handler := &cityRouter{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	//only staff
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.AccessAuth))
		r.Use(jwtauth.Authenticator(conf.AccessAuth))
		r.Use(middleware.StaffPermission)

		r.Post("/", handler.createCity)
		r.With(middleware.Pagination).Get("/", handler.listCities)
		r.Put("/{id}", handler.updateCity)
		r.Delete("/{id}", handler.deleteCity)
	})

	//public
	r.With(middleware.Pagination).Get("/active", handler.listActiveCities)

	return r
}

// Handler's
func (h *cityRouter) createCity(w http.ResponseWriter, r *http.Request) {
	var input cityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	id, err := h.queries.CreateCity(
		r.Context(),
		db.CreateCityParams{
			NameEn:   input.NameEn,
			NameAr:   input.NameAr,
			IsActive: *input.IsActive,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JsonResponseWriter(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *cityRouter) listCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get limit and offset from parsed ctx and filter query
	page := ctx.Value("pagination").(*middleware.Paginator)
	query := r.URL.Query().Get("q")

	var err error
	var cities []db.City
	if query != "" {
		arg := db.FilterAllCitiesParams{
			Limit:  page.Limit,
			Offset: page.Offset,
			Query:  query,
		}
		cities, err = h.queries.FilterAllCities(ctx, arg)
	} else {
		arg := db.ListAllCitiesParams{
			Limit:  page.Limit,
			Offset: page.Offset,
		}
		cities, err = h.queries.ListAllCities(ctx, arg)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]cityResponse, len(cities))
	for i, city := range cities {
		response[i] = cityResponse{
			ID:       city.ID,
			NameEn:   city.NameEn,
			NameAr:   city.NameAr,
			IsActive: city.IsActive,
		}
	}

	totalCount, err := h.queries.CitiesCount(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityRouter) listActiveCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//TODO: try to make ListActiveCityParams AND ListAllCitiesParams one struct
	page := ctx.Value("pagination").(*middleware.Paginator)
	cities, err := h.queries.ListActiveCity(ctx, db.ListActiveCityParams{page.Limit, page.Offset})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestLang := r.Header.Get("Accept-Language")
	response := make([]activeCityResponse, len(cities))
	for i, city := range cities {
		var name string
		if requestLang == "ar" {
			name = city.NameAr
		} else {
			name = city.NameEn
		}

		response[i] = activeCityResponse{
			ID:   city.ID,
			Name: name,
		}
	}

	// get total cities count
	totalCount, err := h.queries.ActiveCityCount(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityRouter) updateCity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input cityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	city, err := h.queries.UpdateCity(
		r.Context(),
		db.UpdateCityParams{
			ID:       id,
			NameEn:   input.NameEn,
			NameAr:   input.NameAr,
			IsActive: *input.IsActive,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "City not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := cityResponse{
		ID:       city.ID,
		NameEn:   city.NameEn,
		NameAr:   city.NameAr,
		IsActive: city.IsActive,
	}

	util.JsonResponseWriter(w, http.StatusOK, response)
}

func (h *cityRouter) deleteCity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.queries.DeleteCity(ctx, id); err != nil {
		http.Error(w, "Failed to delete city", http.StatusInternalServerError)
		return
	}

	util.JsonResponseWriter(w, http.StatusNoContent, nil)
}
