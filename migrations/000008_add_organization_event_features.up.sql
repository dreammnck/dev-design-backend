-- Publish status for events (org submits → admin approves/rejects)
CREATE TYPE event_publish_status AS ENUM ('pending', 'approved', 'rejected');

-- Add organization ownership + publish workflow columns to events
ALTER TABLE events
    ADD COLUMN organization_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN publish_status  event_publish_status NOT NULL DEFAULT 'pending',
    ADD COLUMN reject_reason   TEXT,
    ADD COLUMN published_at    TIMESTAMP WITH TIME ZONE;

-- Add seats per zone / capacity concept for org seat management
-- (seat_type is a free-form label: e.g. "VIP", "General", "Zone A")
ALTER TABLE seats
    ADD COLUMN seat_type VARCHAR(100);

-- Payout requests: org requests money transfer after event sales
CREATE TYPE payout_status AS ENUM ('requested', 'processing', 'completed', 'rejected');

CREATE TABLE payouts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id         UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    amount           INT NOT NULL,
    status           payout_status NOT NULL DEFAULT 'requested',
    bank_account     VARCHAR(255),
    bank_name        VARCHAR(100),
    reject_reason    TEXT,
    requested_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at     TIMESTAMP WITH TIME ZONE,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
