package auth

import "testing"

func TestHasPermission(t *testing.T) {
	if !HasPermission("booking_manager", "booking.confirm") {
		t.Fatal("booking manager must be able to confirm a booking")
	}
	if HasPermission("guest", "booking.update") {
		t.Fatal("guest must not be able to manage bookings")
	}
	if !HasPermission("super_admin", "roles.manage") {
		t.Fatal("super admin must have every permission")
	}
}
