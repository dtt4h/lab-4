package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"lab-4/backend/internal/models"
)

type IncidentCommentsRepository interface {
	Create(ctx context.Context, comment *models.IncidentComment) error
	ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentComment, error)
	Update(ctx context.Context, comment *models.IncidentComment) error
	Delete(ctx context.Context, id string) error
}
type incidentCommentsRepository struct{ db *pgxpool.Pool }

func NewIncidentCommentsRepository(db *pgxpool.Pool) IncidentCommentsRepository {
	return &incidentCommentsRepository{db: db}
}
func (r *incidentCommentsRepository) Create(ctx context.Context, comment *models.IncidentComment) error {
	err := r.db.QueryRow(ctx, `INSERT INTO incident_comments (incident_id, author_id, body) VALUES ($1,$2,$3) RETURNING id, created_at, updated_at`, comment.IncidentID, comment.AuthorID, comment.Body).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	return normalizeError(err)
}
func (r *incidentCommentsRepository) ListByIncident(ctx context.Context, incidentID string) ([]models.IncidentComment, error) {
	rows, err := r.db.Query(ctx, `SELECT id, incident_id, author_id, body, created_at, updated_at FROM incident_comments WHERE incident_id = $1 ORDER BY created_at`, incidentID)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	comments := make([]models.IncidentComment, 0)
	for rows.Next() {
		var comment models.IncidentComment
		if err := rows.Scan(&comment.ID, &comment.IncidentID, &comment.AuthorID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
func (r *incidentCommentsRepository) Update(ctx context.Context, comment *models.IncidentComment) error {
	err := r.db.QueryRow(ctx, `UPDATE incident_comments SET body = $2, updated_at = NOW() WHERE id = $1 RETURNING updated_at`, comment.ID, comment.Body).Scan(&comment.UpdatedAt)
	return normalizeError(err)
}
func (r *incidentCommentsRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM incident_comments WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
