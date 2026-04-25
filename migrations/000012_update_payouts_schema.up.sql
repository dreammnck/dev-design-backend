-- Add account_name column and make event_id nullable
ALTER TABLE payouts 
    ADD COLUMN account_name VARCHAR(255),
    ALTER COLUMN event_id DROP NOT NULL;
