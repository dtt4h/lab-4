package handler

import (
	"encoding/json"
	"lab-4/backend/internal/middleware"
	"lab-4/backend/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct{ users service.UsersService }

func NewAdminHandler(users service.UsersService) *AdminHandler { return &AdminHandler{users} }
func (h *AdminHandler) SetRole(w http.ResponseWriter, r *http.Request) {
	if !h.isSuper(r) {
		writeForbidden(w)
		return
	}
	var body struct {
		Role string `json:"role"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.users.UpdateRole(r.Context(), chi.URLParam(r, "id"), body.Role); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.isSuper(r) {
		writeForbidden(w)
		return
	}
	items, err := h.users.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}
func (h *AdminHandler) isSuper(r *http.Request) bool {
	id, _ := middleware.UserID(r.Context())
	u, err := h.users.GetByID(r.Context(), id)
	return err == nil && u.Role == "super_admin"
}
