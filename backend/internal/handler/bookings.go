package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"lab-4/backend/internal/mailer"
	"lab-4/backend/internal/middleware"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/service"
)

type BookingsHandler struct {
	bookings service.BookingsService
	rooms    service.RoomsService
	users    service.UsersService
	mailer   mailer.BookingMailer
}

func NewBookingsHandler(b service.BookingsService, r service.RoomsService, u service.UsersService, m mailer.BookingMailer) *BookingsHandler {
	return &BookingsHandler{bookings: b, rooms: r, users: u, mailer: m}
}
func (h *BookingsHandler) Rooms(w http.ResponseWriter, r *http.Request) {
	items, err := h.rooms.ListAvailable(r.Context(), r.URL.Query().Get("hotel_id"), r.URL.Query().Get("check_in"), r.URL.Query().Get("check_out"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *BookingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	var b models.Booking
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	b.GuestID = userID
	if err := h.bookings.Create(r.Context(), &b); err != nil {
		writeError(w, err)
		return
	}
	h.sendReceipt(b)
	writeJSON(w, 201, map[string]any{"data": b})
}
func (h *BookingsHandler) PublicCreate(w http.ResponseWriter, r *http.Request) {
	var b models.Booking
	if json.NewDecoder(r.Body).Decode(&b) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.bookings.Create(r.Context(), &b); err != nil {
		writeError(w, err)
		return
	}
	h.sendReceipt(b)
	writeJSON(w, http.StatusCreated, map[string]any{"data": b})
}

func (h *BookingsHandler) sendReceipt(booking models.Booking) {
	if h.mailer == nil || booking.GuestEmail == "" {
		return
	}
	room, err := h.rooms.GetByID(context.Background(), booking.RoomID)
	if err != nil {
		log.Printf("booking receipt skipped: room lookup failed for booking %s", booking.ID)
		return
	}
	receipt := mailer.BookingReceipt{To: booking.GuestEmail, GuestName: booking.GuestName, RoomNumber: room.Number, RoomType: room.RoomTypeName, CheckIn: booking.CheckIn, CheckOut: booking.CheckOut, GuestsCount: booking.GuestsCount, TotalPrice: booking.TotalPrice}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.mailer.SendBookingReceipt(ctx, receipt); err != nil {
			log.Printf("booking receipt delivery failed for booking %s", booking.ID)
		}
	}()
}
func (h *BookingsHandler) RoomTypes(w http.ResponseWriter, r *http.Request) {
	items, err := h.rooms.ListTypes(r.Context(), r.URL.Query().Get("hotel_id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}
func (h *BookingsHandler) CreateRoomType(w http.ResponseWriter, r *http.Request) {
	var x models.RoomType
	if json.NewDecoder(r.Body).Decode(&x) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.rooms.CreateType(r.Context(), &x); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": x})
}
func (h *BookingsHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var x models.Room
	if json.NewDecoder(r.Body).Decode(&x) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.rooms.Create(r.Context(), &x); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": x})
}
func (h *BookingsHandler) My(w http.ResponseWriter, r *http.Request) {
	id, _ := middleware.UserID(r.Context())
	items, err := h.bookings.ListByGuest(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *BookingsHandler) Pay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b, err := h.bookings.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	userID, _ := middleware.UserID(r.Context())
	if b.GuestID != userID {
		writeForbidden(w)
		return
	}
	if err := h.bookings.UpdateStatus(r.Context(), id, "paid"); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": map[string]string{"status": "paid"}})
}
func (h *BookingsHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	if !h.canManage(r) {
		writeForbidden(w)
		return
	}
	items, err := h.bookings.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": items})
}
func (h *BookingsHandler) AdminStatus(w http.ResponseWriter, r *http.Request) {
	if !h.canManage(r) {
		writeForbidden(w)
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		writeError(w, service.ErrValidation)
		return
	}
	if err := h.bookings.UpdateStatus(r.Context(), chi.URLParam(r, "id"), body.Status); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(204)
}
func (h *BookingsHandler) canManage(r *http.Request) bool {
	id, _ := middleware.UserID(r.Context())
	u, err := h.users.GetByID(r.Context(), id)
	return err == nil && (u.Role == "super_admin" || u.Role == "booking_manager")
}
