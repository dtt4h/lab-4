ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('super_admin', 'booking_manager', 'security_manager', 'hotel_manager', 'room_manager', 'guest'));
