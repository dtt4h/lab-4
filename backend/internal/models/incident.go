package models

import "time"

type Incident struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	HotelID     string     `json:"hotel_id"`
	LocationID  *string    `json:"location_id"`
	CategoryID  *string    `json:"category_id"`
	CreatedBy   string     `json:"created_by"`
	OccurredAt  time.Time  `json:"occurred_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type IncidentFilter struct {
	HotelID  *string
	Status   *string
	Priority *string
	Limit    int
	Offset   int
}
