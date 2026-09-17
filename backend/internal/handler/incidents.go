package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"lab-4/backend/internal/handler/dto"
	"lab-4/backend/internal/middleware"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/service"
)

type IncidentsHandler struct{ incidents service.IncidentsService }

func NewIncidentsHandler(incidents service.IncidentsService) *IncidentsHandler {
	return &IncidentsHandler{incidents: incidents}
}

func (h *IncidentsHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.IncidentFilter{HotelID: optional(r.URL.Query().Get("hotel_id")), Status: optional(r.URL.Query().Get("status")), Priority: optional(r.URL.Query().Get("priority")), Limit: queryInt(r, "limit", 20), Offset: (queryInt(r, "page", 1) - 1) * queryInt(r, "limit", 20)}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	items, err := h.incidents.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	page, limit := queryInt(r, "page", 1), filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items, "pagination": map[string]int{"page": page, "limit": limit, "total": len(items)}})
}
func (h *IncidentsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	incident, err := h.incidents.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": incident})
}
func (h *IncidentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": map[string]string{"code": "unauthorized", "message": "Authentication required"}})
		return
	}
	var req dto.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	incident := incidentFromRequest(req)
	incident.CreatedBy = userID
	if err := h.incidents.Create(r.Context(), incident); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": incident})
}
func (h *IncidentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	current, err := h.incidents.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	if !owns(r, current.CreatedBy) {
		writeForbidden(w)
		return
	}
	var req dto.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	updated := incidentFromRequest(req)
	updated.ID = id
	updated.CreatedBy = current.CreatedBy
	if err := h.incidents.Update(r.Context(), updated); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": updated})
}
func (h *IncidentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	incident, err := h.incidents.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	if !owns(r, incident.CreatedBy) {
		writeForbidden(w)
		return
	}
	if err := h.incidents.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func incidentFromRequest(req dto.IncidentRequest) *models.Incident {
	return &models.Incident{Title: strings.TrimSpace(req.Title), Description: strings.TrimSpace(req.Description), Priority: req.Priority, Status: req.Status, HotelID: req.HotelID, LocationID: req.LocationID, CategoryID: req.CategoryID, OccurredAt: req.OccurredAt, ResolvedAt: req.ResolvedAt}
}
func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
func queryInt(r *http.Request, name string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value == 0 {
		return fallback
	}
	return value
}
func owns(r *http.Request, ownerID string) bool {
	id, ok := middleware.UserID(r.Context())
	return ok && id == ownerID
}
func writeForbidden(w http.ResponseWriter) {
	writeJSON(w, http.StatusForbidden, map[string]any{"error": map[string]string{"code": "forbidden", "message": "You do not have permission to modify this incident"}})
}
