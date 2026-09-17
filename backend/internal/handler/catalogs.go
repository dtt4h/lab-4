package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/service"
)

type CatalogsHandler struct {
	hotels     service.HotelsService
	locations  service.LocationsService
	categories service.IncidentCategoriesService
}

func NewCatalogsHandler(h service.HotelsService, l service.LocationsService, c service.IncidentCategoriesService) *CatalogsHandler {
	return &CatalogsHandler{hotels: h, locations: l, categories: c}
}
func (h *CatalogsHandler) Hotels(w http.ResponseWriter, r *http.Request) {
	items, err := h.hotels.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *CatalogsHandler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var v models.Hotel
	if json.NewDecoder(r.Body).Decode(&v) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	v.Name = strings.TrimSpace(v.Name)
	v.Address = strings.TrimSpace(v.Address)
	if err := h.hotels.Create(r.Context(), &v); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"data": v})
}
func (h *CatalogsHandler) Categories(w http.ResponseWriter, r *http.Request) {
	items, err := h.categories.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *CatalogsHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var v models.IncidentCategory
	if json.NewDecoder(r.Body).Decode(&v) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	v.Name = strings.TrimSpace(v.Name)
	if err := h.categories.Create(r.Context(), &v); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"data": v})
}
func (h *CatalogsHandler) Locations(w http.ResponseWriter, r *http.Request) {
	hotelID := r.URL.Query().Get("hotel_id")
	items, err := h.locations.ListByHotel(r.Context(), hotelID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *CatalogsHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var v models.Location
	if json.NewDecoder(r.Body).Decode(&v) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	v.Floor = strings.TrimSpace(v.Floor)
	v.Room = strings.TrimSpace(v.Room)
	v.AreaName = strings.TrimSpace(v.AreaName)
	if err := h.locations.Create(r.Context(), &v); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"data": v})
}
func (h *CatalogsHandler) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	if err := h.hotels.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(204)
}
func (h *CatalogsHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if err := h.categories.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(204)
}
func (h *CatalogsHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	if err := h.locations.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(204)
}
