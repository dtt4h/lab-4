package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type LocationsService interface {
	Create(context.Context, *models.Location) error
	GetByID(context.Context, string) (*models.Location, error)
	ListByHotel(context.Context, string) ([]models.Location, error)
	Update(context.Context, *models.Location) error
	Delete(context.Context, string) error
}
type locationsService struct {
	locations repository.LocationsRepository
}

func NewLocationsService(locations repository.LocationsRepository) LocationsService {
	return &locationsService{locations: locations}
}
func (s *locationsService) Create(ctx context.Context, v *models.Location) error {
	if err := required(v.HotelID, v.Floor); err != nil {
		return err
	}
	return s.locations.Create(ctx, v)
}
func (s *locationsService) GetByID(ctx context.Context, id string) (*models.Location, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.locations.GetByID(ctx, id)
}
func (s *locationsService) ListByHotel(ctx context.Context, hotelID string) ([]models.Location, error) {
	if err := required(hotelID); err != nil {
		return nil, err
	}
	return s.locations.ListByHotel(ctx, hotelID)
}
func (s *locationsService) Update(ctx context.Context, v *models.Location) error {
	if err := required(v.ID, v.HotelID, v.Floor); err != nil {
		return err
	}
	return s.locations.Update(ctx, v)
}
func (s *locationsService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.locations.Delete(ctx, id)
}
