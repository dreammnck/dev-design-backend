-- Creation of UUID extension (if not already exists)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE seat_status AS ENUM ('available', 'reserved', 'sold');

CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'cancelled');

CREATE TYPE payment_status AS ENUM ('success', 'failed');

CREATE TYPE payment_method AS ENUM ('credit_card', 'qr_code', 'bank_transfer', 'true_money');

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    latitude FLOAT,
    longitude FLOAT,
    city VARCHAR(100),
    state_province VARCHAR(100),
    country VARCHAR(100),
    post_code VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(512),
    location_id UUID REFERENCES locations (id) ON DELETE SET NULL,
    event_date DATE NOT NULL,
    event_time TIME,
    price INT DEFAULT 0,
    is_banner BOOLEAN DEFAULT FALSE,
    is_recommend BOOLEAN DEFAULT FALSE,
    is_coming_soon BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    seat_number VARCHAR(50) NOT NULL,
    status seat_status DEFAULT 'available',
    price INT,
    reserved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    event_id UUID NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    total_amount INT NOT NULL,
    status booking_status DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE booking_seats (
    booking_id UUID NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    seat_id UUID NOT NULL REFERENCES seats (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (booking_id, seat_id)
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    booking_id UUID NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    transaction_id VARCHAR(255),
    amount INT NOT NULL,
    payment_method payment_method,
    status payment_status DEFAULT 'failed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE payment_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    payment_id UUID NOT NULL REFERENCES payments (id) ON DELETE CASCADE,
    detail_key VARCHAR(50) NOT NULL,
    detail_value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);