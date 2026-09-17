package dto

type EmployeeRequest struct {
	UserID   *string `json:"user_id"`
	HotelID  string  `json:"hotel_id"`
	FullName string  `json:"full_name"`
	Position string  `json:"position"`
}
