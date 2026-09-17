package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type IncidentCommentsService interface {
	Create(ctx context.Context, comment *models.IncidentComment) error
	ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentComment, error)
	Update(ctx context.Context, comment *models.IncidentComment) error
	Delete(ctx context.Context, id string) error
}
type incidentCommentsService struct {
	comments repository.IncidentCommentsRepository
}

func NewIncidentCommentsService(comments repository.IncidentCommentsRepository) IncidentCommentsService {
	return &incidentCommentsService{comments: comments}
}
func (s *incidentCommentsService) Create(ctx context.Context, v *models.IncidentComment) error {
	if err := required(v.IncidentID, v.AuthorID, v.Body); err != nil {
		return err
	}
	return s.comments.Create(ctx, v)
}
func (s *incidentCommentsService) ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentComment, error) {
	if err := required(incidentID); err != nil {
		return nil, err
	}
	return s.comments.ListByIncident(ctx, incidentID)
}
func (s *incidentCommentsService) Update(ctx context.Context, v *models.IncidentComment) error {
	if err := required(v.ID, v.Body); err != nil {
		return err
	}
	return s.comments.Update(ctx, v)
}
func (s *incidentCommentsService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.comments.Delete(ctx, id)
}
