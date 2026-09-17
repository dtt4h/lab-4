package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"lab-4/backend/internal/handler/dto"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/service"
)

type EmployeesHandler struct{ employees service.EmployeesService }

func NewEmployeesHandler(employees service.EmployeesService) *EmployeesHandler {
	return &EmployeesHandler{employees: employees}
}
func (h *EmployeesHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.employees.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}
func (h *EmployeesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	employee, err := h.employees.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": employee})
}
func (h *EmployeesHandler) Create(w http.ResponseWriter, r *http.Request) {
	employee, err := decodeEmployee(r)
	if err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.employees.Create(r.Context(), employee); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": employee})
}
func (h *EmployeesHandler) Update(w http.ResponseWriter, r *http.Request) {
	employee, err := decodeEmployee(r)
	if err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	employee.ID = chi.URLParam(r, "id")
	if err := h.employees.Update(r.Context(), employee); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": employee})
}
func (h *EmployeesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.employees.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func decodeEmployee(r *http.Request) (*models.Employee, error) {
	var req dto.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &models.Employee{UserID: req.UserID, HotelID: req.HotelID, FullName: strings.TrimSpace(req.FullName), Position: strings.TrimSpace(req.Position)}, nil
}
