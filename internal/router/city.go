package router

import (
	"encoding/json"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type cityHandler struct {
	conf    *config.Setting
	queries *db.Queries
}

func newCityRouter(conf *config.Setting, queries *db.Queries) *chi.Mux {
	router := chi.NewRouter()
	handler := &cityHandler{
		queries: queries,
		conf:    conf,
	}

	router.Post("/", handler.createCity)
	router.Get("/", handler.listCities)
	router.Get("/active", handler.listActiveCities)
	router.Put("/{id}", handler.updateCity)
	router.Delete("/{id}", handler.deleteCity)

	return router
}

func (h *cityHandler) createCity(w http.ResponseWriter, r *http.Request) {

	type CityInput struct {
		NameEn   string `json:"name_en" validate:"required"`
		NameAr   string `json:"name_ar" validate:"required"`
		IsActive *bool  `json:"is_active,omitempty"`
	}

	var input CityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	if err := validate.Struct(input); err != nil {
		ts := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			ts[err.Field()] = err.Tag()
		}

		util.JsonResponseWriter(w, http.StatusBadRequest, ts)
		return
	}

	if input.IsActive == nil {
		defaultState := true
		input.IsActive = &defaultState
	}

	params := db.CreateCityParams{
		NameEn:   input.NameEn,
		NameAr:   input.NameAr,
		IsActive: *input.IsActive,
	}

	id, err := h.queries.CreateCity(r.Context(), params)
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.JsonResponseWriter(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *cityHandler) listCities(w http.ResponseWriter, r *http.Request) {
	var l, o int64 = 10, 0

	if limitQueryParam, ok := r.URL.Query()["limit"]; ok && len(limitQueryParam[0]) > 0 {
		limitInt, err := strconv.ParseInt(limitQueryParam[0], 10, 64)
		if err != nil {
			util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
			return
		}
		l = limitInt
	}

	if offsetQueryParam, ok := r.URL.Query()["offset"]; ok && len(offsetQueryParam[0]) > 0 {
		offsetInt, err := strconv.ParseInt(offsetQueryParam[0], 10, 64)
		if err != nil {
			util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
			return
		}
		o = offsetInt
	}

	totalCount, err := h.queries.CitiesCount(r.Context())
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}

	cities, err := h.queries.ListAllCities(r.Context(), db.ListAllCitiesParams{Limit: l, Offset: o})
	if err != nil {
		util.ErrorResponseWriter(w, http.StatusBadRequest, err.Error())
		return
	}
	type cityResponse struct {
		ID       int64  `json:"id"`
		NameEn   string `json:"name_en"`
		NameAr   string `json:"name_ar"`
		IsActive bool   `json:"is_active"`
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

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityHandler) listActiveCities(w http.ResponseWriter, r *http.Request) {

}

func (h *cityHandler) updateCity(w http.ResponseWriter, r *http.Request) {

}

func (h *cityHandler) deleteCity(w http.ResponseWriter, r *http.Request) {

}
