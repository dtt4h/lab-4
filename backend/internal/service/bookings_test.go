package service

import (
	"context"
	"testing"

	"lab-4/backend/internal/models"
)

type bookingRepositoryStub struct {
	booking *models.Booking
	updated string
}

func (r *bookingRepositoryStub) Create(context.Context, *models.Booking) error { return nil }
func (r *bookingRepositoryStub) GetByID(context.Context, string) (*models.Booking, error) {
	return r.booking, nil
}
func (r *bookingRepositoryStub) ListByGuest(context.Context, string) ([]models.Booking, error) {
	return nil, nil
}
func (r *bookingRepositoryStub) List(context.Context) ([]models.Booking, error) { return nil, nil }
func (r *bookingRepositoryStub) UpdateStatus(_ context.Context, _ string, status string) error {
	r.updated = status
	return nil
}

func TestBookingStatusTransitions(t *testing.T) {
	repo := &bookingRepositoryStub{booking: &models.Booking{ID: "booking-1", Status: "paid"}}
	service := NewBookingsService(repo)
	if err := service.UpdateStatus(context.Background(), "booking-1", "confirmed"); err != nil {
		t.Fatalf("paid booking should be confirmable: %v", err)
	}
	if repo.updated != "confirmed" {
		t.Fatalf("expected confirmed, got %q", repo.updated)
	}
	if err := service.UpdateStatus(context.Background(), "booking-1", "pending"); err != ErrValidation {
		t.Fatalf("expected invalid backwards transition, got %v", err)
	}
}
