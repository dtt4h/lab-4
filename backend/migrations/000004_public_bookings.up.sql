ALTER TABLE bookings ALTER COLUMN guest_id DROP NOT NULL;
ALTER TABLE bookings ADD COLUMN guest_name VARCHAR(120) NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN guest_email VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN guest_phone VARCHAR(30) NOT NULL DEFAULT '';
CREATE INDEX bookings_dates_active_idx ON bookings(room_id, check_in, check_out) WHERE status IN ('pending', 'paid', 'confirmed');
