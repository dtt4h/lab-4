package models

import "time"

type RoomType struct {
	ID            string  `json:"id"`
	HotelID       string  `json:"hotel_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Capacity      int     `json:"capacity"`
	PricePerNight float64 `json:"price_per_night"`
}
type Room struct {
	ID            string  `json:"id"`
	HotelID       string  `json:"hotel_id"`
	RoomTypeID    string  `json:"room_type_id"`
	Number        string  `json:"number"`
	Floor         string  `json:"floor"`
	Status        string  `json:"status"`
	RoomTypeName  string  `json:"room_type_name"`
	Capacity      int     `json:"capacity"`
	PricePerNight float64 `json:"price_per_night"`
}
type Booking struct {
	ID          string    `json:"id"`
	GuestID     string    `json:"guest_id,omitempty"`
	GuestName   string    `json:"guest_name"`
	GuestEmail  string    `json:"guest_email"`
	GuestPhone  string    `json:"guest_phone"`
	RoomID      string    `json:"room_id"`
	CheckIn     string    `json:"check_in"`
	CheckOut    string    `json:"check_out"`
	GuestsCount int       `json:"guests_count"`
	Status      string    `json:"status"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
