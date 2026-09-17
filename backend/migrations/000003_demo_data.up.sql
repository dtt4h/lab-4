INSERT INTO hotels (name, address) VALUES ('Cosmos Hotel', 'ул. Центральная, 1') ON CONFLICT (name) DO NOTHING;
INSERT INTO room_types (hotel_id, name, description, capacity, price_per_night)
SELECT id, 'Стандарт', 'Уютный номер для отдыха', 2, 4200 FROM hotels WHERE name='Cosmos Hotel'
ON CONFLICT (hotel_id, name) DO NOTHING;
INSERT INTO room_types (hotel_id, name, description, capacity, price_per_night)
SELECT id, 'Семейный', 'Просторный номер для семьи', 4, 7200 FROM hotels WHERE name='Cosmos Hotel'
ON CONFLICT (hotel_id, name) DO NOTHING;
INSERT INTO rooms (hotel_id, room_type_id, number, floor)
SELECT h.id, rt.id, '101', '1' FROM hotels h JOIN room_types rt ON rt.hotel_id=h.id AND rt.name='Стандарт' WHERE h.name='Cosmos Hotel'
ON CONFLICT (hotel_id, number) DO NOTHING;
INSERT INTO rooms (hotel_id, room_type_id, number, floor)
SELECT h.id, rt.id, '102', '1' FROM hotels h JOIN room_types rt ON rt.hotel_id=h.id AND rt.name='Семейный' WHERE h.name='Cosmos Hotel'
ON CONFLICT (hotel_id, number) DO NOTHING;
