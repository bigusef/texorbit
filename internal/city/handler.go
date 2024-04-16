package city

import (
	"encoding/json"
	"errors"
	db "github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/bigusef/texorbit/pkg/middleware"
	"github.com/bigusef/texorbit/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type cityHandler struct {
	conf     *config.Setting
	queries  *db.Queries
	validate *validator.Validate
}

func (h *cityHandler) createCity(w http.ResponseWriter, r *http.Request) {
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

func (h *cityHandler) listCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get limit and offset from parsed ctx and filter query
	page := ctx.Value("pagination").(*middleware.Paginator)
	query := r.URL.Query().Get("q")

	var err error
	var cities []db.City
	if query != "" {
		arg := db.FilterCitiesParams{
			Limit:  page.Limit,
			Offset: page.Offset,
			Query:  query,
		}
		cities, err = h.queries.FilterCities(ctx, arg)
	} else {
		arg := db.AllCitiesParams{
			Limit:  page.Limit,
			Offset: page.Offset,
		}
		cities, err = h.queries.AllCities(ctx, arg)
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

func (h *cityHandler) listActiveCities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//TODO: try to make ListActiveCityParams AND ListAllCitiesParams one struct
	page := ctx.Value("pagination").(*middleware.Paginator)
	cities, err := h.queries.ActiveCities(ctx, db.ActiveCitiesParams{page.Limit, page.Offset})
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
	totalCount, err := h.queries.ActiveCitiesCount(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JsonListResponseWriter(w, http.StatusOK, response, totalCount)
}

func (h *cityHandler) updateCity(w http.ResponseWriter, r *http.Request) {
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

func (h *cityHandler) deleteCity(w http.ResponseWriter, r *http.Request) {
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
