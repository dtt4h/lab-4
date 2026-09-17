package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lab-4/backend/internal/models"
)

type RoomsRepository interface {
	ListAvailable(context.Context, string, string, string) ([]models.Room, error)
	GetByID(context.Context, string) (*models.Room, error)
	ListTypes(context.Context, string) ([]models.RoomType, error)
	CreateType(context.Context, *models.RoomType) error
	Create(context.Context, *models.Room) error
}

func (r *roomsRepository) ListTypes(ctx context.Context, hotelID string) ([]models.RoomType, error) {
	rows, err := r.db.Query(ctx, `SELECT id,hotel_id,name,description,capacity,price_per_night FROM room_types WHERE hotel_id=$1 ORDER BY price_per_night`, hotelID)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	items := []models.RoomType{}
	for rows.Next() {
		var x models.RoomType
		if err := rows.Scan(&x.ID, &x.HotelID, &x.Name, &x.Description, &x.Capacity, &x.PricePerNight); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}
func (r *roomsRepository) CreateType(ctx context.Context, x *models.RoomType) error {
	return normalizeError(r.db.QueryRow(ctx, `INSERT INTO room_types (hotel_id,name,description,capacity,price_per_night) VALUES ($1,$2,$3,$4,$5) RETURNING id`, x.HotelID, x.Name, x.Description, x.Capacity, x.PricePerNight).Scan(&x.ID))
}
func (r *roomsRepository) Create(ctx context.Context, x *models.Room) error {
	return normalizeError(r.db.QueryRow(ctx, `INSERT INTO rooms (hotel_id,room_type_id,number,floor) VALUES ($1,$2,$3,$4) RETURNING id,status`, x.HotelID, x.RoomTypeID, x.Number, x.Floor).Scan(&x.ID, &x.Status))
}

type BookingsRepository interface {
	Create(context.Context, *models.Booking) error
	GetByID(context.Context, string) (*models.Booking, error)
	ListByGuest(context.Context, string) ([]models.Booking, error)
	List(context.Context) ([]models.Booking, error)
	UpdateStatus(context.Context, string, string) error
}
type roomsRepository struct{ db *pgxpool.Pool }
type bookingsRepository struct{ db *pgxpool.Pool }

func NewRoomsRepository(db *pgxpool.Pool) RoomsRepository       { return &roomsRepository{db: db} }
func NewBookingsRepository(db *pgxpool.Pool) BookingsRepository { return &bookingsRepository{db: db} }

func (r *roomsRepository) GetByID(ctx context.Context, id string) (*models.Room, error) {
	room := new(models.Room)
	err := r.db.QueryRow(ctx, `SELECT r.id,r.hotel_id,r.room_type_id,r.number,r.floor,r.status,rt.name,rt.capacity,rt.price_per_night FROM rooms r JOIN room_types rt ON rt.id=r.room_type_id WHERE r.id=$1`, id).Scan(&room.ID, &room.HotelID, &room.RoomTypeID, &room.Number, &room.Floor, &room.Status, &room.RoomTypeName, &room.Capacity, &room.PricePerNight)
	if err != nil {
		return nil, normalizeError(err)
	}
	return room, nil
}

func (r *roomsRepository) ListAvailable(ctx context.Context, hotelID, checkIn, checkOut string) ([]models.Room, error) {
	const query = `SELECT r.id,r.hotel_id,r.room_type_id,r.number,r.floor,r.status,rt.name,rt.capacity,rt.price_per_night
		FROM rooms r JOIN room_types rt ON rt.id=r.room_type_id
		WHERE r.hotel_id=$1 AND r.status='available'
		AND NOT EXISTS (SELECT 1 FROM bookings b WHERE b.room_id=r.id AND b.status IN ('pending','paid','confirmed') AND b.check_in < $3::date AND b.check_out > $2::date)
		ORDER BY rt.price_per_night, r.number`
	rows, err := r.db.Query(ctx, query, hotelID, checkIn, checkOut)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	items := []models.Room{}
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.RoomTypeID, &room.Number, &room.Floor, &room.Status, &room.RoomTypeName, &room.Capacity, &room.PricePerNight); err != nil {
			return nil, err
		}
		items = append(items, room)
	}
	return items, rows.Err()
}

func (r *bookingsRepository) Create(ctx context.Context, b *models.Booking) error {
	const query = `INSERT INTO bookings (guest_id,guest_name,guest_email,guest_phone,room_id,check_in,check_out,guests_count,total_price)
		SELECT NULLIF($1,'')::uuid,$2,$3,$4,r.id,$6,$7,$8,rt.price_per_night*($7::date-$6::date)
		FROM rooms r JOIN room_types rt ON rt.id=r.room_type_id
		WHERE r.id=$5 AND r.status='available' AND rt.capacity >= $8
		AND NOT EXISTS (SELECT 1 FROM bookings b WHERE b.room_id=$5 AND b.status IN ('pending','paid','confirmed') AND b.check_in < $7 AND b.check_out > $6)
		RETURNING id,status,total_price,created_at,updated_at`
	err := r.db.QueryRow(ctx, query, b.GuestID, b.GuestName, b.GuestEmail, b.GuestPhone, b.RoomID, b.CheckIn, b.CheckOut, b.GuestsCount).Scan(&b.ID, &b.Status, &b.TotalPrice, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrConflict
	}
	return normalizeError(err)
}
func (r *bookingsRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	b := new(models.Booking)
	err := r.db.QueryRow(ctx, `SELECT id,COALESCE(guest_id::text,''),guest_name,guest_email,guest_phone,room_id,check_in::text,check_out::text,guests_count,status,total_price,created_at,updated_at FROM bookings WHERE id=$1`, id).Scan(&b.ID, &b.GuestID, &b.GuestName, &b.GuestEmail, &b.GuestPhone, &b.RoomID, &b.CheckIn, &b.CheckOut, &b.GuestsCount, &b.Status, &b.TotalPrice, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, normalizeError(err)
	}
	return b, nil
}
func (r *bookingsRepository) ListByGuest(ctx context.Context, id string) ([]models.Booking, error) {
	return r.list(ctx, `WHERE guest_id=$1 ORDER BY created_at DESC`, id)
}
func (r *bookingsRepository) List(ctx context.Context) ([]models.Booking, error) {
	return r.list(ctx, `ORDER BY created_at DESC`)
}
func (r *bookingsRepository) list(ctx context.Context, suffix string, args ...any) ([]models.Booking, error) {
	rows, err := r.db.Query(ctx, `SELECT id,COALESCE(guest_id::text,''),guest_name,guest_email,guest_phone,room_id,check_in::text,check_out::text,guests_count,status,total_price,created_at,updated_at FROM bookings `+suffix, args...)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	items := []models.Booking{}
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.GuestID, &b.GuestName, &b.GuestEmail, &b.GuestPhone, &b.RoomID, &b.CheckIn, &b.CheckOut, &b.GuestsCount, &b.Status, &b.TotalPrice, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}
func (r *bookingsRepository) UpdateStatus(ctx context.Context, id, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE bookings SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	if err != nil {
		return normalizeError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
