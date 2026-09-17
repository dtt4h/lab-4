package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"lab-4/backend/internal/models"
)

type IncidentsRepository interface {
	Create(context.Context, *models.Incident) error
	GetByID(context.Context, string) (*models.Incident, error)
	List(context.Context, models.IncidentFilter) ([]models.Incident, error)
	Update(context.Context, *models.Incident) error
	Delete(context.Context, string) error
}
type incidentsRepository struct{ db *pgxpool.Pool }

func NewIncidentsRepository(db *pgxpool.Pool) IncidentsRepository {
	return &incidentsRepository{db: db}
}
func (r *incidentsRepository) Create(ctx context.Context, incident *models.Incident) error {
	const query = `INSERT INTO incidents (title, description, priority, status, hotel_id, location_id, category_id, created_by, occurred_at, resolved_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, incident.Title, incident.Description, incident.Priority, incident.Status, incident.HotelID, incident.LocationID, incident.CategoryID, incident.CreatedBy, incident.OccurredAt, incident.ResolvedAt).Scan(&incident.ID, &incident.CreatedAt, &incident.UpdatedAt)
	return normalizeError(err)
}
func (r *incidentsRepository) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	incident := new(models.Incident)
	err := r.db.QueryRow(ctx, incidentSelect+` WHERE id = $1`, id).Scan(incidentFields(incident)...)
	if err != nil {
		return nil, normalizeError(err)
	}
	return incident, nil
}
func (r *incidentsRepository) List(ctx context.Context, filter models.IncidentFilter) ([]models.Incident, error) {
	query := incidentSelect + ` WHERE ($1::uuid IS NULL OR hotel_id = $1) AND ($2::text IS NULL OR status = $2) AND ($3::text IS NULL OR priority = $3) ORDER BY occurred_at DESC LIMIT $4 OFFSET $5`
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.Query(ctx, query, filter.HotelID, filter.Status, filter.Priority, limit, offset)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	incidents := make([]models.Incident, 0)
	for rows.Next() {
		var incident models.Incident
		if err := rows.Scan(incidentFields(&incident)...); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incidents, nil
}
func (r *incidentsRepository) Update(ctx context.Context, incident *models.Incident) error {
	const query = `UPDATE incidents SET title=$2, description=$3, priority=$4, status=$5, hotel_id=$6, location_id=$7, category_id=$8, occurred_at=$9, resolved_at=$10, updated_at=NOW() WHERE id=$1 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, incident.ID, incident.Title, incident.Description, incident.Priority, incident.Status, incident.HotelID, incident.LocationID, incident.CategoryID, incident.OccurredAt, incident.ResolvedAt).Scan(&incident.UpdatedAt)
	return normalizeError(err)
}
func (r *incidentsRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM incidents WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

const incidentSelect = `SELECT id, title, description, priority, status, hotel_id, location_id, category_id, created_by, occurred_at, resolved_at, created_at, updated_at FROM incidents`

func incidentFields(incident *models.Incident) []any {
	return []any{&incident.ID, &incident.Title, &incident.Description, &incident.Priority, &incident.Status, &incident.HotelID, &incident.LocationID, &incident.CategoryID, &incident.CreatedBy, &incident.OccurredAt, &incident.ResolvedAt, &incident.CreatedAt, &incident.UpdatedAt}
}
