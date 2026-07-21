# Poll app

`docker compose up --build` runs PostgreSQL and the Go server. Copy `.env.example` to `.env` first.

API: `POST /api/polls` with `{ "title": "...", "questions": [{"text":"...", "options":["A","B"]}] }`; `GET /api/polls/{id}`; `POST /api/polls/{id}/vote` with `{ "answers": [{"question_id":"...", "option_id":"..."}] }`; and `GET /api/polls/{id}/results` (optionally `?admin_token=...`).
