# Agent Guide

## Overview

This repository is a Go/chi + PostgreSQL poll API with a Vue/Vite SPA. The server creates multi-question UUID polls, accepts one answer per question, and serves public/admin results.

## Ownership and layers

- `cmd/server`: process startup, environment configuration, migration execution, static/SPA serving.
- `internal/domain`: transport-independent poll models.
- `internal/repository`: storage interfaces and PostgreSQL implementation. SQL belongs in `internal/repository/postgres`.
- `internal/service`: validation and business rules.
- `internal/handler`: HTTP/JSON, status mapping, and cookies; never couple handlers to PostgreSQL.
- `migrations`: SQL schema files executed by the server.
- `frontend`: Vue routes, views, and browser API client.

Keep API JSON names in `snake_case` and preserve the documented endpoint shapes, status codes, error object, admin-token behavior, and `voted_{poll_id}` cookie contract.

## Verification commands

```sh
go test ./...
go build ./cmd/server
cd frontend && npm ci && npm run build
cd .. && docker compose config
docker compose up --build
docker build .
```

Run backend locally with PostgreSQL available and `DATABASE_URL` set. Run the frontend with `cd frontend && npm run dev`.

## Constraints and style

- Use pgx/pgxpool for PostgreSQL; do not add another database driver.
- Keep handlers independent of concrete PostgreSQL packages.
- Validate at the service boundary and keep multi-table poll creation transactional.
- Add schema changes as idempotent migration SQL; do not edit production data manually.
- Prefer small, idiomatic Go changes and focused tests for non-trivial validation.
- Do not commit generated `frontend/dist`, binaries, logs, Docker build output, or dependency caches.
- Do not change frontend routes or API payloads without updating the compatibility documentation and all callers.
