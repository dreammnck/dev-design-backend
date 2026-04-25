ALTER TABLE events ADD COLUMN user_id UUID;
-- Optional: if you want to link it to users
-- ALTER TABLE events ADD CONSTRAINT fk_events_user_id FOREIGN KEY (user_id) REFERENCES users(id);
