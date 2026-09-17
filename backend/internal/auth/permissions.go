package auth

// Permissions are evaluated on the server. The UI receives them only to
// present an appropriate experience; it is never an authorization boundary.
var rolePermissions = map[string][]string{
	"super_admin":      {"*"},
	"booking_manager":  {"booking.read", "booking.update", "booking.confirm", "booking.cancel"},
	"room_manager":     {"room.read", "room.create"},
	"hotel_manager":    {"hotel.read", "hotel.manage", "employee.manage", "catalog.manage"},
	"security_manager": {"incident.read", "incident.manage", "hotel.read", "catalog.read"},
	"guest":            {"booking.create", "booking.read.own", "booking.pay"},
}

func PermissionsForRole(role string) []string {
	permissions := rolePermissions[role]
	return append([]string(nil), permissions...)
}

func HasPermission(role, permission string) bool {
	for _, item := range rolePermissions[role] {
		if item == "*" || item == permission {
			return true
		}
	}
	return false
}
