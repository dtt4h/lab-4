package models

import "time"

type IncidentComment struct {
	ID         string
	IncidentID string
	AuthorID   string
	Body       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
