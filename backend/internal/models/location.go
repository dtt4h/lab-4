package models

type Location struct {
	ID       string `json:"id"`
	HotelID  string `json:"hotel_id"`
	Floor    string `json:"floor"`
	Room     string `json:"room"`
	AreaName string `json:"area_name"`
}
