CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS polls (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), title text NOT NULL, admin_token text NOT NULL UNIQUE, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS questions (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), poll_id uuid NOT NULL REFERENCES polls(id) ON DELETE CASCADE, text text NOT NULL, position int NOT NULL);
CREATE TABLE IF NOT EXISTS options (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), question_id uuid NOT NULL REFERENCES questions(id) ON DELETE CASCADE, text text NOT NULL, position int NOT NULL, UNIQUE(question_id,text));
CREATE TABLE IF NOT EXISTS votes (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), poll_id uuid NOT NULL REFERENCES polls(id) ON DELETE CASCADE, question_id uuid NOT NULL REFERENCES questions(id) ON DELETE CASCADE, option_id uuid NOT NULL REFERENCES options(id) ON DELETE CASCADE);
CREATE INDEX IF NOT EXISTS votes_option_id_idx ON votes(option_id);
