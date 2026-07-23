DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'polls' AND column_name = 'results_visible_to_voter')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'polls' AND column_name = 'results_visible') THEN
        ALTER TABLE polls RENAME COLUMN results_visible_to_voter TO results_visible;
    ELSIF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'polls' AND column_name = 'results_visible_to_voter') THEN
        ALTER TABLE polls DROP COLUMN results_visible_to_voter;
    END IF;
END $$;
ALTER TABLE polls DROP COLUMN IF EXISTS anonymous;
ALTER TABLE polls ADD COLUMN IF NOT EXISTS results_visible boolean NOT NULL DEFAULT true;
ALTER TABLE polls ALTER COLUMN results_visible SET DEFAULT true;
UPDATE polls SET results_visible = true WHERE results_visible IS NULL;
ALTER TABLE polls ALTER COLUMN results_visible SET NOT NULL;
