package service

import (
	"context"
	"strings"
	"time"

	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type BookingsService interface {
	Create(context.Context, *models.Booking) error
	GetByID(context.Context, string) (*models.Booking, error)
	ListByGuest(context.Context, string) ([]models.Booking, error)
	List(context.Context) ([]models.Booking, error)
	UpdateStatus(context.Context, string, string) error
}
type RoomsService interface {
	ListAvailable(context.Context, string, string, string) ([]models.Room, error)
	GetByID(context.Context, string) (*models.Room, error)
	ListTypes(context.Context, string) ([]models.RoomType, error)
	CreateType(context.Context, *models.RoomType) error
	Create(context.Context, *models.Room) error
}

func (s *roomsService) GetByID(ctx context.Context, id string) (*models.Room, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *roomsService) ListTypes(ctx context.Context, id string) ([]models.RoomType, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.repo.ListTypes(ctx, id)
}
func (s *roomsService) CreateType(ctx context.Context, x *models.RoomType) error {
	x.Name = strings.TrimSpace(x.Name)
	x.Description = strings.TrimSpace(x.Description)
	if err := required(x.HotelID, x.Name); err != nil || x.Capacity < 1 || x.PricePerNight < 0 {
		return ErrValidation
	}
	return s.repo.CreateType(ctx, x)
}
func (s *roomsService) Create(ctx context.Context, x *models.Room) error {
	if err := required(x.HotelID, x.RoomTypeID, x.Number, x.Floor); err != nil {
		return err
	}
	return s.repo.Create(ctx, x)
}

type bookingsService struct{ repo repository.BookingsRepository }
type roomsService struct{ repo repository.RoomsRepository }

func NewBookingsService(repo repository.BookingsRepository) BookingsService {
	return &bookingsService{repo}
}
func NewRoomsService(repo repository.RoomsRepository) RoomsService { return &roomsService{repo} }
func (s *roomsService) ListAvailable(ctx context.Context, hotelID, checkIn, checkOut string) ([]models.Room, error) {
	if err := required(hotelID, checkIn, checkOut); err != nil {
		return nil, err
	}
	in, err := time.Parse("2006-01-02", checkIn)
	if err != nil {
		return nil, ErrValidation
	}
	out, err := time.Parse("2006-01-02", checkOut)
	if err != nil || !out.After(in) {
		return nil, ErrValidation
	}
	return s.repo.ListAvailable(ctx, hotelID, checkIn, checkOut)
}
func (s *bookingsService) Create(ctx context.Context, b *models.Booking) error {
	if err := required(b.RoomID, b.CheckIn, b.CheckOut); err != nil {
		return err
	}
	if b.GuestID == "" && (strings.TrimSpace(b.GuestName) == "" || !strings.Contains(b.GuestEmail, "@")) {
		return ErrValidation
	}
	if b.GuestsCount < 1 {
		return ErrValidation
	}
	in, err := time.Parse("2006-01-02", b.CheckIn)
	if err != nil {
		return ErrValidation
	}
	out, err := time.Parse("2006-01-02", b.CheckOut)
	if err != nil || !out.After(in) {
		return ErrValidation
	}
	return s.repo.Create(ctx, b)
}
func (s *bookingsService) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *bookingsService) ListByGuest(ctx context.Context, id string) ([]models.Booking, error) {
	return s.repo.ListByGuest(ctx, id)
}
func (s *bookingsService) List(ctx context.Context) ([]models.Booking, error) {
	return s.repo.List(ctx)
}
func (s *bookingsService) UpdateStatus(ctx context.Context, id, status string) error {
	if status != "pending" && status != "paid" && status != "confirmed" && status != "cancelled" {
		return ErrValidation
	}
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	allowed := map[string]map[string]bool{
		"pending":   {"paid": true, "cancelled": true},
		"paid":      {"confirmed": true, "cancelled": true},
		"confirmed": {"cancelled": true},
	}
	if !allowed[booking.Status][status] {
		return ErrValidation
	}
	return s.repo.UpdateStatus(ctx, id, status)
}
