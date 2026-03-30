-- Seed data for locations
INSERT INTO locations (id, name, latitude, longitude, city, state_province, country, post_code, is_active)
VALUES (
    '110e8400-e29b-41d4-a716-446655440001',
    'Bangkok Arena',
    13.7563,
    100.5018,
    'Bangkok',
    'Bangkok',
    'Thailand',
    '10100',
    TRUE
);

-- Seed data for events
INSERT INTO events (id, title, description, image_url, location_id, event_date, event_time, price, is_banner, is_recommend)
VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    'Music Festival 2026',
    'A big music festival with many artists.',
    'https://example.com/images/music-festival-2026.jpg',
    '110e8400-e29b-41d4-a716-446655440001',
    '2026-04-10',
    '18:00:00',
    1500,
    TRUE,
    TRUE
);

-- Seed data for seats
INSERT INTO seats (event_id, seat_number, status)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'A1', 'available'),
    ('550e8400-e29b-41d4-a716-446655440000', 'A2', 'available'),
    ('550e8400-e29b-41d4-a716-446655440000', 'B1', 'available');
