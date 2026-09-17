package repository

import (
	"context"
	"lab-4/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HotelsRepository interface {
	Create(context.Context, *models.Hotel) error
	GetByID(context.Context, string) (*models.Hotel, error)
	List(context.Context) ([]models.Hotel, error)
	Update(context.Context, *models.Hotel) error
	Delete(context.Context, string) error
}

type hotelsRepository struct{ db *pgxpool.Pool }

func NewHotelsRepository(db *pgxpool.Pool) HotelsRepository {
	return &hotelsRepository{db: db}

}

func (r *hotelsRepository) Create(ctx context.Context, hotel *models.Hotel) error {
	err := r.db.QueryRow(ctx, `INSERT INTO hotels (name, address) VALUES ($1, $2) RETURNING id`, hotel.Name, hotel.Address).Scan(&hotel.ID)
	return normalizeError(err)
}

func (r *hotelsRepository) GetByID(ctx context.Context, id string) (*models.Hotel, error) {
	hotel := new(models.Hotel)
	err := r.db.QueryRow(ctx, `SELECT id, name, address FROM hotels WHERE id = $1`, id).Scan(&hotel.ID, &hotel.Name, &hotel.Address)
	if err != nil {
		return nil, normalizeError(err)
	}
	return hotel, nil
}

func (r *hotelsRepository) List(ctx context.Context) ([]models.Hotel, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, address FROM hotels ORDER BY name`)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	hotels := make([]models.Hotel, 0)
	for rows.Next() {
		var hotel models.Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (r *hotelsRepository) Update(ctx context.Context, hotel *models.Hotel) error {
	result, err := r.db.Exec(ctx, `UPDATE hotels SET name = $2, address = $3 WHERE id = $1`, hotel.ID, hotel.Name, hotel.Address)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *hotelsRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM hotels WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
