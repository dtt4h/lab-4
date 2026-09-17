package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type IncidentsService interface {
	Create(context.Context, *models.Incident) error
	GetByID(context.Context, string) (*models.Incident, error)
	List(context.Context, models.IncidentFilter) ([]models.Incident, error)
	Update(context.Context, *models.Incident) error
	Delete(context.Context, string) error
}
type incidentsService struct {
	incidents repository.IncidentsRepository
}

func NewIncidentsService(incidents repository.IncidentsRepository) IncidentsService {
	return &incidentsService{incidents: incidents}
}
func (s *incidentsService) Create(ctx context.Context, v *models.Incident) error {
	if err := required(v.Title, v.Priority, v.Status, v.HotelID, v.CreatedBy); err != nil {
		return err
	}
	return s.incidents.Create(ctx, v)
}
func (s *incidentsService) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.incidents.GetByID(ctx, id)
}
func (s *incidentsService) List(ctx context.Context, filter models.IncidentFilter) ([]models.Incident, error) {
	return s.incidents.List(ctx, filter)
}
func (s *incidentsService) Update(ctx context.Context, v *models.Incident) error {
	if err := required(v.ID, v.Title, v.Priority, v.Status, v.HotelID); err != nil {
		return err
	}
	return s.incidents.Update(ctx, v)
}
func (s *incidentsService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.incidents.Delete(ctx, id)
}
