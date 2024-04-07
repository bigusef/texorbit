package router

import (
	"encoding/json"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

// Schema data DTOs
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

// Services
func getPaginationParams(r *http.Request) (int64, int64, error) {
	var l, o int64 = 10, 0

	if param, ok := r.URL.Query()["limit"]; ok && len(param[0]) > 0 {
		limitInt, err := strconv.ParseInt(param[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		l = limitInt
	}

	if param, ok := r.URL.Query()["offset"]; ok && len(param[0]) > 0 {
		offsetInt, err := strconv.ParseInt(param[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		o = offsetInt
	}
	return l, o, nil
}

func NewCityRouter(conf *config.Setting, queries *db.Queries, validate *validator.Validate) *chi.Mux {
	router := chi.NewRouter()
	handler := &cityRouter{
		queries:  queries,
		conf:     conf,
		validate: validate,
	}

	router.Post("/", handler.createCity)
	router.Get("/", handler.listCities)
	router.Get("/active", handler.listActiveCities)
	router.Put("/{id}", handler.updateCity)
	router.Delete("/{id}", handler.deleteCity)

	return router
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
	limit, offset, err := getPaginationParams(r)
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
	limit, offset, err := getPaginationParams(r)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
	}

	//TODO: try to make ListActiveCityParams AND ListAllCitiesParams one struct
	cities, err := h.queries.ListActiveCity(r.Context(), db.ListActiveCityParams{
		Limit:  offset,
		Offset: limit,
	})
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]activeCityResponse, len(cities))
	for i, city := range cities {
		response[i] = activeCityResponse{
			ID:   city.ID,
			Name: city.NameEn, // TODO: return name based on city
		}
	}

	totalCount, err := h.queries.CitiesCount(r.Context())
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
		// TODO: cover case row not found and issue in update process
		util.ErrorResponseWriter(w, http.StatusNotFound, err.Error())
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

	err = h.queries.DeleteCity(r.Context(), id)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.JsonResponseWriter(w, http.StatusNoContent, nil)
}
