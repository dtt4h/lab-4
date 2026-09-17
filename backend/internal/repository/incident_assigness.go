package repository

import (
	"context"
	"lab-4/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentAssignmentsRepository interface {
	Assign(ctx context.Context, assignment *models.IncidentAssignment) error
	Unassign(ctx context.Context, incidentID, employeeID string) error
	ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentAssignment, error)
}
type incidentAssignmentsRepository struct{ db *pgxpool.Pool }

func NewIncidentAssignmentsRepository(db *pgxpool.Pool) IncidentAssignmentsRepository {
	return &incidentAssignmentsRepository{db: db}
}

func (r *incidentAssignmentsRepository) Assign(ctx context.Context, assignment *models.IncidentAssignment) error {
	err := r.db.QueryRow(ctx, `INSERT INTO incident_assignments (incident_id, employee_id, assigned_by) VALUES ($1, $2, $3) RETURNING assigned_at`, assignment.IncidentID, assignment.EmployeeID, assignment.AssignedBy).Scan(&assignment.AssignedAt)
	return normalizeError(err)
}
func (r *incidentAssignmentsRepository) Unassign(ctx context.Context, incidentID, employeeID string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM incident_assignments WHERE incident_id = $1 AND employee_id = $2`, incidentID, employeeID)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *incidentAssignmentsRepository) ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentAssignment, error) {
	rows, err := r.db.Query(ctx, `SELECT incident_id, employee_id, assigned_at, assigned_by FROM incident_assignments WHERE incident_id = $1 ORDER BY assigned_at`, incidentID)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	assignments := make([]models.IncidentAssignment, 0)
	for rows.Next() {
		var assignment models.IncidentAssignment
		if err := rows.Scan(&assignment.IncidentID, &assignment.EmployeeID, &assignment.AssignedAt, &assignment.AssignedBy); err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assignments, nil
}
