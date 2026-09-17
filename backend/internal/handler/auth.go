package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"lab-4/backend/internal/auth"
	"lab-4/backend/internal/handler/dto"
	"lab-4/backend/internal/middleware"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/service"
)

type AuthHandler struct {
	users      service.UsersService
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthHandler(users service.UsersService, jwtSecret string, accessTTL time.Duration) *AuthHandler {
	return &AuthHandler{users: users, jwtSecret: jwtSecret, accessTTL: accessTTL, refreshTTL: 7 * 24 * time.Hour}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if len(req.Password) < 8 {
		writeError(w, service.ErrValidation)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, err)
		return
	}
	user, err := h.users.Register(r.Context(), models.CreateUserParams{Username: strings.TrimSpace(req.Username), Email: strings.TrimSpace(req.Email), PasswordHash: hash})
	if err != nil {
		writeError(w, err)
		return
	}
	h.writeAuthResponse(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, service.ErrValidation)
		return
	}
	user, err := h.users.GetByLogin(r.Context(), strings.TrimSpace(req.Login))
	if err != nil || auth.VerifyPassword(user.PasswordHash, req.Password) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": map[string]string{"code": "invalid_credentials", "message": "Invalid login or password"}})
		return
	}
	if err := h.users.UpdateLastLogin(r.Context(), user.ID); err != nil {
		writeError(w, err)
		return
	}
	h.writeAuthResponse(w, http.StatusOK, user)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": userResponse(user)})
}
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": map[string]string{"code": "unauthorized", "message": "Authentication required"}})
		return
	}
	id, err := auth.ParseRefreshToken(cookie.Value, h.jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": map[string]string{"code": "unauthorized", "message": "Authentication required"}})
		return
	}
	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	h.writeAuthResponse(w, http.StatusOK, user)
}
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/api/v1/auth", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	w.WriteHeader(http.StatusNoContent)
}
func (h *AuthHandler) writeAuthResponse(w http.ResponseWriter, status int, user *models.User) {
	token, err := auth.NewAccessToken(user.ID, h.jwtSecret, h.accessTTL)
	if err != nil {
		writeError(w, err)
		return
	}
	refresh, err := auth.NewRefreshToken(user.ID, h.jwtSecret, h.refreshTTL)
	if err != nil {
		writeError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: refresh, Path: "/api/v1/auth", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: int(h.refreshTTL.Seconds())})
	writeJSON(w, status, map[string]any{"data": dto.AuthResponse{AccessToken: token, User: userResponse(user)}})
}
func userResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role, Permissions: auth.PermissionsForRole(user.Role)}
}
