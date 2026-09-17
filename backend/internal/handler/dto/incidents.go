package dto

import "time"

type IncidentRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	HotelID     string     `json:"hotel_id"`
	LocationID  *string    `json:"location_id"`
	CategoryID  *string    `json:"category_id"`
	OccurredAt  time.Time  `json:"occurred_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}
