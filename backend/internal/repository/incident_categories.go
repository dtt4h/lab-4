package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"lab-4/backend/internal/models"
)

type IncidentCategoriesRepository interface {
	Create(context.Context, *models.IncidentCategory) error
	GetByID(context.Context, string) (*models.IncidentCategory, error)
	List(context.Context) ([]models.IncidentCategory, error)
	Update(context.Context, *models.IncidentCategory) error
	Delete(context.Context, string) error
}
type incidentCategoriesRepository struct{ db *pgxpool.Pool }

func NewIncidentCategoriesRepository(db *pgxpool.Pool) IncidentCategoriesRepository {
	return &incidentCategoriesRepository{db: db}
}
func (r *incidentCategoriesRepository) Create(ctx context.Context, category *models.IncidentCategory) error {
	err := r.db.QueryRow(ctx, `INSERT INTO incident_categories (name) VALUES ($1) RETURNING id`, category.Name).Scan(&category.ID)
	return normalizeError(err)
}
func (r *incidentCategoriesRepository) GetByID(ctx context.Context, id string) (*models.IncidentCategory, error) {
	category := new(models.IncidentCategory)
	err := r.db.QueryRow(ctx, `SELECT id, name FROM incident_categories WHERE id = $1`, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, normalizeError(err)
	}
	return category, nil
}
func (r *incidentCategoriesRepository) List(ctx context.Context) ([]models.IncidentCategory, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM incident_categories ORDER BY name`)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	categories := make([]models.IncidentCategory, 0)
	for rows.Next() {
		var category models.IncidentCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
func (r *incidentCategoriesRepository) Update(ctx context.Context, category *models.IncidentCategory) error {
	result, err := r.db.Exec(ctx, `UPDATE incident_categories SET name = $2 WHERE id = $1`, category.ID, category.Name)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *incidentCategoriesRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM incident_categories WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
