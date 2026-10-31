# Roshandel – resume site

What it is: public home page + per-user resume at /<slug></slug>, editable in a panel; roles user / admin / owner.

## Run locally

make dev / make web (ports, .env.local, seed admin via ADMIN_EMAIL/ADMIN_PASSWORD)

## Deploy

Push to main → CI tests, builds arm64 image → GHCR → Watchtower on the server pulls it.
Server: ~/app, docker-compose.prod.yml, .env (list every variable), Caddy for TLS, DOMAIN.
Manual: docker compose -f docker-compose.prod.yml pull app && up -d app
Rollback: pin image: tag to a commit sha.

## Data

Migrations: goose, auto-run at start; make migration name=...
Backup: pg_dump command; restore command.
Seeding a resume: how, and that it targets the first user.

## Layout

(keep the existing tree)
