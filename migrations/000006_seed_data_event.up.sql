
INSERT INTO locations (id, name, latitude, longitude, city, state_province, country, post_code, is_active)
VALUES
('110e8400-e29b-41d4-a716-446655440002','Chiang Mai Hall',18.7883,98.9853,'Chiang Mai','Chiang Mai','Thailand','50000',TRUE),
('110e8400-e29b-41d4-a716-446655440003','Phuket Beach Stage',7.8804,98.3923,'Phuket','Phuket','Thailand','83000',TRUE);


INSERT INTO events (id, title, description, image_url, location_id, event_date, event_time, price, is_banner, is_recommend, is_coming_soon)
VALUES
('550e8400-e29b-41d4-a716-446655440040','EDM Festival Night','Join an unforgettable EDM festival with world-class DJs, immersive lighting, and high-energy performances throughout the night. Experience music, visuals, and crowd energy like never before.','https://images.unsplash.com/photo-1507874457470-272b3c8d8ee2','110e8400-e29b-41d4-a716-446655440002','2026-06-10','20:00:00',800,TRUE,TRUE,FALSE),

('550e8400-e29b-41d4-a716-446655440041','Jazz & Chill Evening','Enjoy a relaxing evening filled with smooth jazz performances in a cozy atmosphere. Perfect for music lovers who want to unwind and enjoy live instrumental music.','https://images.unsplash.com/photo-1500530855697-b586d89ba3ee','110e8400-e29b-41d4-a716-446655440002','2026-06-12','18:30:00',1000,FALSE,FALSE,FALSE),

('550e8400-e29b-41d4-a716-446655440042','Beach Sunset Concert','Experience live acoustic performances by the beach during sunset. Enjoy scenic views, relaxing vibes, and great music with friends and family.','https://images.unsplash.com/photo-1497032205916-ac775f0649ae','110e8400-e29b-41d4-a716-446655440003','2026-06-15','17:00:00',1000,TRUE,TRUE,FALSE),

('550e8400-e29b-41d4-a716-446655440043','Stand-up Comedy Night','Laugh the night away with top comedians delivering hilarious stand-up performances. A fun and entertaining event for everyone.','https://images.unsplash.com/photo-1527224857830-43a7acc85260','110e8400-e29b-41d4-a716-446655440002','2026-06-18','20:00:00',700,FALSE,TRUE,FALSE),

('550e8400-e29b-41d4-a716-446655440044','Tech & AI Conference','Explore the future of technology with expert speakers discussing AI, software, and innovation. Network with professionals and attend insightful sessions.','https://images.unsplash.com/photo-1518770660439-4636190af475','110e8400-e29b-41d4-a716-446655440002','2026-06-20','09:00:00',2000,TRUE,FALSE,TRUE),

('550e8400-e29b-41d4-a716-446655440045','Food Carnival','Taste a variety of street food and international cuisines. Enjoy cooking shows, food tasting, and entertainment activities throughout the day.','https://images.unsplash.com/photo-1504674900247-0877df9cc836','110e8400-e29b-41d4-a716-446655440003','2026-06-22','11:00:00',500,FALSE,TRUE,FALSE),

('550e8400-e29b-41d4-a716-446655440046','Art & Design Expo','Discover creative works from artists and designers. Explore exhibitions, installations, and interactive art displays in a modern setting.','https://images.unsplash.com/photo-1492724441997-5dc865305da7','110e8400-e29b-41d4-a716-446655440002','2026-06-25','10:00:00',600,FALSE,FALSE,TRUE),

('550e8400-e29b-41d4-a716-446655440047','Outdoor Movie Night','Enjoy a relaxing outdoor cinema experience under the stars. Watch popular movies with friends and family in a comfortable open-air setting.','https://images.unsplash.com/photo-1505685296765-3a2736de412f','110e8400-e29b-41d4-a716-446655440003','2026-06-28','19:30:00',400,FALSE,FALSE,FALSE);


INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440040','A1','available',2000),
('550e8400-e29b-41d4-a716-446655440040','A2','available',2000),
('550e8400-e29b-41d4-a716-446655440040','A3','available',2000),
('550e8400-e29b-41d4-a716-446655440040','A4','available',2000),
('550e8400-e29b-41d4-a716-446655440040','A5','available',2000),
('550e8400-e29b-41d4-a716-446655440040','A6','available',1500),
('550e8400-e29b-41d4-a716-446655440040','A7','available',1500),
('550e8400-e29b-41d4-a716-446655440040','A8','available',1500),
('550e8400-e29b-41d4-a716-446655440040','A9','available',1500),
('550e8400-e29b-41d4-a716-446655440040','A10','available',1500),
('550e8400-e29b-41d4-a716-446655440040','B1','available',1000),
('550e8400-e29b-41d4-a716-446655440040','B2','available',1000),
('550e8400-e29b-41d4-a716-446655440040','B3','available',1000),
('550e8400-e29b-41d4-a716-446655440040','B4','available',1000),
('550e8400-e29b-41d4-a716-446655440040','B5','available',1000),
('550e8400-e29b-41d4-a716-446655440040','B6','available',800),
('550e8400-e29b-41d4-a716-446655440040','B7','available',800),
('550e8400-e29b-41d4-a716-446655440040','B8','available',800),
('550e8400-e29b-41d4-a716-446655440040','B9','available',800),
('550e8400-e29b-41d4-a716-446655440040','B10','available',800);




INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440041','A1','available',2500),
('550e8400-e29b-41d4-a716-446655440041','A2','available',2500),
('550e8400-e29b-41d4-a716-446655440041','A3','available',2500),
('550e8400-e29b-41d4-a716-446655440041','A4','available',2500),
('550e8400-e29b-41d4-a716-446655440041','A5','available',2500),
('550e8400-e29b-41d4-a716-446655440041','A6','available',2000),
('550e8400-e29b-41d4-a716-446655440041','A7','available',2000),
('550e8400-e29b-41d4-a716-446655440041','A8','available',2000),
('550e8400-e29b-41d4-a716-446655440041','A9','available',2000),
('550e8400-e29b-41d4-a716-446655440041','A10','available',2000),
('550e8400-e29b-41d4-a716-446655440041','B1','available',1500),
('550e8400-e29b-41d4-a716-446655440041','B2','available',1500),
('550e8400-e29b-41d4-a716-446655440041','B3','available',1500),
('550e8400-e29b-41d4-a716-446655440041','B4','available',1500),
('550e8400-e29b-41d4-a716-446655440041','B5','available',1500),
('550e8400-e29b-41d4-a716-446655440041','B6','available',1000),
('550e8400-e29b-41d4-a716-446655440041','B7','available',1000),
('550e8400-e29b-41d4-a716-446655440041','B8','available',1000),
('550e8400-e29b-41d4-a716-446655440041','B9','available',1000),
('550e8400-e29b-41d4-a716-446655440041','B10','available',1000);


INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440042','A1','available',2500),
('550e8400-e29b-41d4-a716-446655440042','A2','available',2500),
('550e8400-e29b-41d4-a716-446655440042','A3','available',2500),
('550e8400-e29b-41d4-a716-446655440042','A4','available',2500),
('550e8400-e29b-41d4-a716-446655440042','A5','available',2500),
('550e8400-e29b-41d4-a716-446655440042','A6','available',2000),
('550e8400-e29b-41d4-a716-446655440042','A7','available',2000),
('550e8400-e29b-41d4-a716-446655440042','A8','available',2000),
('550e8400-e29b-41d4-a716-446655440042','A9','available',2000),
('550e8400-e29b-41d4-a716-446655440042','A10','available',2000),
('550e8400-e29b-41d4-a716-446655440042','B1','available',1500),
('550e8400-e29b-41d4-a716-446655440042','B2','available',1500),
('550e8400-e29b-41d4-a716-446655440042','B3','available',1500),
('550e8400-e29b-41d4-a716-446655440042','B4','available',1500),
('550e8400-e29b-41d4-a716-446655440042','B5','available',1500),
('550e8400-e29b-41d4-a716-446655440042','B6','available',1000),
('550e8400-e29b-41d4-a716-446655440042','B7','available',1000),
('550e8400-e29b-41d4-a716-446655440042','B8','available',1000),
('550e8400-e29b-41d4-a716-446655440042','B9','available',1000),
('550e8400-e29b-41d4-a716-446655440042','B10','available',1000);

INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440043','A1','available',3000),
('550e8400-e29b-41d4-a716-446655440043','A2','available',3000),
('550e8400-e29b-41d4-a716-446655440043','A3','available',3000),
('550e8400-e29b-41d4-a716-446655440043','A4','available',3000),
('550e8400-e29b-41d4-a716-446655440043','A5','available',3000),
('550e8400-e29b-41d4-a716-446655440043','A6','available',2000),
('550e8400-e29b-41d4-a716-446655440043','A7','available',2000),
('550e8400-e29b-41d4-a716-446655440043','A8','available',2000),
('550e8400-e29b-41d4-a716-446655440043','A9','available',2000),
('550e8400-e29b-41d4-a716-446655440043','A10','available',2000),
('550e8400-e29b-41d4-a716-446655440043','B1','available',1500),
('550e8400-e29b-41d4-a716-446655440043','B2','available',1500),
('550e8400-e29b-41d4-a716-446655440043','B3','available',1500),
('550e8400-e29b-41d4-a716-446655440043','B4','available',1500),
('550e8400-e29b-41d4-a716-446655440043','B5','available',1500),
('550e8400-e29b-41d4-a716-446655440043','B6','available',700),
('550e8400-e29b-41d4-a716-446655440043','B7','available',700),
('550e8400-e29b-41d4-a716-446655440043','B8','available',700),
('550e8400-e29b-41d4-a716-446655440043','B9','available',700),
('550e8400-e29b-41d4-a716-446655440043','B10','available',700);



INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440044','A1','available',4500),
('550e8400-e29b-41d4-a716-446655440044','A2','available',4500),
('550e8400-e29b-41d4-a716-446655440044','A3','available',4500),
('550e8400-e29b-41d4-a716-446655440044','A4','available',4500),
('550e8400-e29b-41d4-a716-446655440044','A5','available',4500),
('550e8400-e29b-41d4-a716-446655440044','A6','available',3000),
('550e8400-e29b-41d4-a716-446655440044','A7','available',3000),
('550e8400-e29b-41d4-a716-446655440044','A8','available',3000),
('550e8400-e29b-41d4-a716-446655440044','A9','available',3000),
('550e8400-e29b-41d4-a716-446655440044','A10','available',3000),
('550e8400-e29b-41d4-a716-446655440044','B1','available',1500),
('550e8400-e29b-41d4-a716-446655440044','B2','available',1500),
('550e8400-e29b-41d4-a716-446655440044','B3','available',1500),
('550e8400-e29b-41d4-a716-446655440044','B4','available',1500),
('550e8400-e29b-41d4-a716-446655440044','B5','available',1500),
('550e8400-e29b-41d4-a716-446655440044','B6','available',2000),
('550e8400-e29b-41d4-a716-446655440044','B7','available',2000),
('550e8400-e29b-41d4-a716-446655440044','B8','available',2000),
('550e8400-e29b-41d4-a716-446655440044','B9','available',2000),
('550e8400-e29b-41d4-a716-446655440044','B10','available',2000);

INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440045','A1','available',1000),
('550e8400-e29b-41d4-a716-446655440045','A2','available',1000),
('550e8400-e29b-41d4-a716-446655440045','A3','available',1000),
('550e8400-e29b-41d4-a716-446655440045','A4','available',1000),
('550e8400-e29b-41d4-a716-446655440045','A5','available',1000),
('550e8400-e29b-41d4-a716-446655440045','A6','available',800),
('550e8400-e29b-41d4-a716-446655440045','A7','available',800),
('550e8400-e29b-41d4-a716-446655440045','A8','available',800),
('550e8400-e29b-41d4-a716-446655440045','A9','available',800),
('550e8400-e29b-41d4-a716-446655440045','A10','available',800),
('550e8400-e29b-41d4-a716-446655440045','B1','available',500),
('550e8400-e29b-41d4-a716-446655440045','B2','available',500),
('550e8400-e29b-41d4-a716-446655440045','B3','available',500),
('550e8400-e29b-41d4-a716-446655440045','B4','available',500),
('550e8400-e29b-41d4-a716-446655440045','B5','available',500),
('550e8400-e29b-41d4-a716-446655440045','B6','available',500),
('550e8400-e29b-41d4-a716-446655440045','B7','available',500),
('550e8400-e29b-41d4-a716-446655440045','B8','available',500),
('550e8400-e29b-41d4-a716-446655440045','B9','available',500),
('550e8400-e29b-41d4-a716-446655440045','B10','available',500);


INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440046','A1','available',1000),
('550e8400-e29b-41d4-a716-446655440046','A2','available',1000),
('550e8400-e29b-41d4-a716-446655440046','A3','available',1000),
('550e8400-e29b-41d4-a716-446655440046','A4','available',1000),
('550e8400-e29b-41d4-a716-446655440046','A5','available',1000),
('550e8400-e29b-41d4-a716-446655440046','A6','available',800),
('550e8400-e29b-41d4-a716-446655440046','A7','available',800),
('550e8400-e29b-41d4-a716-446655440046','A8','available',800),
('550e8400-e29b-41d4-a716-446655440046','A9','available',800),
('550e8400-e29b-41d4-a716-446655440046','A10','available',800),
('550e8400-e29b-41d4-a716-446655440046','B1','available',600),
('550e8400-e29b-41d4-a716-446655440046','B2','available',600),
('550e8400-e29b-41d4-a716-446655440046','B3','available',600),
('550e8400-e29b-41d4-a716-446655440046','B4','available',600),
('550e8400-e29b-41d4-a716-446655440046','B5','available',600),
('550e8400-e29b-41d4-a716-446655440046','B6','available',600),
('550e8400-e29b-41d4-a716-446655440046','B7','available',600),
('550e8400-e29b-41d4-a716-446655440046','B8','available',600),
('550e8400-e29b-41d4-a716-446655440046','B9','available',600),
('550e8400-e29b-41d4-a716-446655440046','B10','available',600);



INSERT INTO seats (event_id, seat_number, status, price) VALUES
('550e8400-e29b-41d4-a716-446655440047','A1','available',1200),
('550e8400-e29b-41d4-a716-446655440047','A2','available',1200),
('550e8400-e29b-41d4-a716-446655440047','A3','available',1200),
('550e8400-e29b-41d4-a716-446655440047','A4','available',1200),
('550e8400-e29b-41d4-a716-446655440047','A5','available',1200),
('550e8400-e29b-41d4-a716-446655440047','A6','available',800),
('550e8400-e29b-41d4-a716-446655440047','A7','available',800),
('550e8400-e29b-41d4-a716-446655440047','A8','available',800),
('550e8400-e29b-41d4-a716-446655440047','A9','available',800),
('550e8400-e29b-41d4-a716-446655440047','A10','available',800),
('550e8400-e29b-41d4-a716-446655440047','B1','available',400),
('550e8400-e29b-41d4-a716-446655440047','B2','available',400),
('550e8400-e29b-41d4-a716-446655440047','B3','available',400),
('550e8400-e29b-41d4-a716-446655440047','B4','available',400),
('550e8400-e29b-41d4-a716-446655440047','B5','available',400),
('550e8400-e29b-41d4-a716-446655440047','B6','available',400),
('550e8400-e29b-41d4-a716-446655440047','B7','available',400),
('550e8400-e29b-41d4-a716-446655440047','B8','available',400),
('550e8400-e29b-41d4-a716-446655440047','B9','available',400),
('550e8400-e29b-41d4-a716-446655440047','B10','available',400);
