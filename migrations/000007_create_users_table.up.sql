CREATE TYPE user_role AS ENUM ('admin', 'organization', 'customer');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'customer',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Mock users (password = "password123" hashed with bcrypt cost 10)
-- Use a placeholder hash — will be replaced by the seed script or app boot
INSERT INTO users (username, email, password_hash, role) VALUES
    ('admin',        'admin@example.com',        '$2a$10$YcXiqyxYgYHe7kphjA9sSuih6yYGehuQoA2dUber0U/Z7bIVfvLgS', 'admin'),
    ('org_user',     'org@example.com',           '$2a$10$YcXiqyxYgYHe7kphjA9sSuih6yYGehuQoA2dUber0U/Z7bIVfvLgS', 'organization'),
    ('customer_one', 'customer1@example.com',     '$2a$10$YcXiqyxYgYHe7kphjA9sSuih6yYGehuQoA2dUber0U/Z7bIVfvLgS', 'customer');
