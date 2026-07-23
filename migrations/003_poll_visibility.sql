ALTER TABLE polls ADD COLUMN IF NOT EXISTS anonymous boolean NOT NULL DEFAULT false;
ALTER TABLE polls ADD COLUMN IF NOT EXISTS results_visible_to_voter boolean NOT NULL DEFAULT false;
