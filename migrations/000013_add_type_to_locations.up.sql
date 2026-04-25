-- Create location_type enum if it doesn't exist
DO $$ BEGIN
    CREATE TYPE location_type AS ENUM ('a', 'b', 'c', 'd');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add column if not exists, and set default
ALTER TABLE locations ADD COLUMN IF NOT EXISTS type location_type NOT NULL DEFAULT 'a';

-- In case column existed but was different, ensure default is set
ALTER TABLE locations ALTER COLUMN type SET DEFAULT 'a';
ALTER TABLE locations ALTER COLUMN type SET NOT NULL;
