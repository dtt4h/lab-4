package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type HotelsService interface {
	Create(context.Context, *models.Hotel) error
	GetByID(context.Context, string) (*models.Hotel, error)
	List(context.Context) ([]models.Hotel, error)
	Update(context.Context, *models.Hotel) error
	Delete(context.Context, string) error
}
type hotelsService struct{ hotels repository.HotelsRepository }

func NewHotelsService(hotels repository.HotelsRepository) HotelsService {
	return &hotelsService{hotels: hotels}
}
func (s *hotelsService) Create(ctx context.Context, hotel *models.Hotel) error {
	if err := required(hotel.Name, hotel.Address); err != nil {
		return err
	}
	return s.hotels.Create(ctx, hotel)
}
func (s *hotelsService) GetByID(ctx context.Context, id string) (*models.Hotel, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.hotels.GetByID(ctx, id)
}
func (s *hotelsService) List(ctx context.Context) ([]models.Hotel, error) { return s.hotels.List(ctx) }
func (s *hotelsService) Update(ctx context.Context, hotel *models.Hotel) error {
	if err := required(hotel.ID, hotel.Name, hotel.Address); err != nil {
		return err
	}
	return s.hotels.Update(ctx, hotel)
}
func (s *hotelsService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.hotels.Delete(ctx, id)
}
