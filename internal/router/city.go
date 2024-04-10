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

		r.Post("/", handler.createCity)
		r.Get("/", handler.listCities)
		r.Put("/{id}", handler.updateCity)
		r.Delete("/{id}", handler.deleteCity)
	})

	//public
	r.Get("/active", handler.listActiveCities)

	return r
}

// Handler's
func (h *cityRouter) createCity(w http.ResponseWriter, r *http.Request) {
	var input cityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
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
		util.ErrorResponseWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.JsonResponseWriter(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *cityRouter) listCities(w http.ResponseWriter, r *http.Request) {
	// TODO: convert to middleware
	// 	https://github.com/go-chi/chi/blob/c1f2a7a12e66b0ca0faf57adad98bf82c767df1f/_examples/rest/main.go#L248
	limit, offset, err := util.GetPaginationParams(r)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
	}

	cities, err := h.queries.ListAllCities(r.Context(), db.ListAllCitiesParams{Limit: limit, Offset: offset})
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
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

	totalCount, err := h.queries.CitiesCount(r.Context())
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityRouter) listActiveCities(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := util.GetPaginationParams(r)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
	}

	//TODO: try to make ListActiveCityParams AND ListAllCitiesParams one struct
	cities, err := h.queries.ListActiveCity(r.Context(), db.ListActiveCityParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
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
	totalCount, err := h.queries.ActiveCityCount(r.Context())
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityRouter) updateCity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	var input cityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
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
			util.ErrorResponseWriter(w, http.StatusNotFound, "City not found")
			return
		}

		util.ErrorResponseWriter(w, http.StatusInternalServerError, err.Error())
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	//TODO: handle if id not exists
	err = h.queries.DeleteCity(r.Context(), id)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.JsonResponseWriter(w, http.StatusNoContent, nil)
}
