package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"lab-4/backend/internal/models"
)

type LocationsRepository interface {
	Create(context.Context, *models.Location) error
	GetByID(context.Context, string) (*models.Location, error)
	ListByHotel(context.Context, string) ([]models.Location, error)
	Update(context.Context, *models.Location) error
	Delete(context.Context, string) error
}
type locationsRepository struct{ db *pgxpool.Pool }

func NewLocationsRepository(db *pgxpool.Pool) LocationsRepository {
	return &locationsRepository{db: db}
}
func (r *locationsRepository) Create(ctx context.Context, location *models.Location) error {
	err := r.db.QueryRow(ctx, `INSERT INTO locations (hotel_id, floor, room, area_name) VALUES ($1, $2, $3, $4) RETURNING id`, location.HotelID, location.Floor, location.Room, location.AreaName).Scan(&location.ID)
	return normalizeError(err)
}
func (r *locationsRepository) GetByID(ctx context.Context, id string) (*models.Location, error) {
	location := new(models.Location)
	err := r.db.QueryRow(ctx, `SELECT id, hotel_id, floor, room, area_name FROM locations WHERE id = $1`, id).Scan(&location.ID, &location.HotelID, &location.Floor, &location.Room, &location.AreaName)
	if err != nil {
		return nil, normalizeError(err)
	}
	return location, nil
}
func (r *locationsRepository) ListByHotel(ctx context.Context, hotelID string) ([]models.Location, error) {
	rows, err := r.db.Query(ctx, `SELECT id, hotel_id, floor, room, area_name FROM locations WHERE hotel_id = $1 ORDER BY floor, room, area_name`, hotelID)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	locations := make([]models.Location, 0)
	for rows.Next() {
		var location models.Location
		if err := rows.Scan(&location.ID, &location.HotelID, &location.Floor, &location.Room, &location.AreaName); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
}
func (r *locationsRepository) Update(ctx context.Context, location *models.Location) error {
	result, err := r.db.Exec(ctx, `UPDATE locations SET hotel_id = $2, floor = $3, room = $4, area_name = $5 WHERE id = $1`, location.ID, location.HotelID, location.Floor, location.Room, location.AreaName)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *locationsRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM locations WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
