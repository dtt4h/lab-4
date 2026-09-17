package models

import "time"

type IncidentAssignment struct {
	IncidentID string
	EmployeeID string
	AssignedAt time.Time
	AssignedBy string
}
