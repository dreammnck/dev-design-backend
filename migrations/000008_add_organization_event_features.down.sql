DROP TABLE IF EXISTS payouts;
DROP TYPE IF EXISTS payout_status;

ALTER TABLE seats DROP COLUMN IF EXISTS seat_type;

ALTER TABLE events
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS reject_reason,
    DROP COLUMN IF EXISTS publish_status,
    DROP COLUMN IF EXISTS organization_id;

DROP TYPE IF EXISTS event_publish_status;
