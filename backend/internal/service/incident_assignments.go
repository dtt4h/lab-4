package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type IncidentAssignmentsService interface {
	Assign(ctx context.Context, assignment *models.IncidentAssignment) error
	Unassign(ctx context.Context, incidentID, employeeID string) error
	ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentAssignment, error)
}
type incidentAssignmentsService struct {
	assignments repository.IncidentAssignmentsRepository
}

func NewIncidentAssignmentsService(assignments repository.IncidentAssignmentsRepository) IncidentAssignmentsService {
	return &incidentAssignmentsService{assignments: assignments}
}
func (s *incidentAssignmentsService) Assign(ctx context.Context, v *models.IncidentAssignment) error {
	if err := required(v.IncidentID, v.EmployeeID, v.AssignedBy); err != nil {
		return err
	}
	return s.assignments.Assign(ctx, v)
}
func (s *incidentAssignmentsService) Unassign(ctx context.Context, incidentID, employeeID string) error {
	if err := required(incidentID, employeeID); err != nil {
		return err
	}
	return s.assignments.Unassign(ctx, incidentID, employeeID)
}
func (s *incidentAssignmentsService) ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentAssignment, error) {
	if err := required(incidentID); err != nil {
		return nil, err
	}
	return s.assignments.ListByIncident(ctx, incidentID)
}
