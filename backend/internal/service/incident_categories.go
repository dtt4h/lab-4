package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type IncidentCategoriesService interface {
	Create(context.Context, *models.IncidentCategory) error
	GetByID(context.Context, string) (*models.IncidentCategory, error)
	List(context.Context) ([]models.IncidentCategory, error)
	Update(context.Context, *models.IncidentCategory) error
	Delete(context.Context, string) error
}
type incidentCategoriesService struct {
	categories repository.IncidentCategoriesRepository
}

func NewIncidentCategoriesService(categories repository.IncidentCategoriesRepository) IncidentCategoriesService {
	return &incidentCategoriesService{categories: categories}
}
func (s *incidentCategoriesService) Create(ctx context.Context, v *models.IncidentCategory) error {
	if err := required(v.Name); err != nil {
		return err
	}
	return s.categories.Create(ctx, v)
}
func (s *incidentCategoriesService) GetByID(ctx context.Context, id string) (*models.IncidentCategory, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.categories.GetByID(ctx, id)
}
func (s *incidentCategoriesService) List(ctx context.Context) ([]models.IncidentCategory, error) {
	return s.categories.List(ctx)
}
func (s *incidentCategoriesService) Update(ctx context.Context, v *models.IncidentCategory) error {
	if err := required(v.ID, v.Name); err != nil {
		return err
	}
	return s.categories.Update(ctx, v)
}
func (s *incidentCategoriesService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.categories.Delete(ctx, id)
}
