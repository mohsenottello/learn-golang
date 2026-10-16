dev:          ## infra in Docker, Go on host with hot reload
	docker compose up -d db mailpit
	set -a; . ./.env.local; set +a; air

web:          ## React dev server
	cd web && npm run dev

up:           ## full stack in Docker
	docker compose up --build -d

down:
	docker compose down

test:
	go test ./...

migration:    ## make migration name=add_something
	@test -n "$(name)" || (echo "usage: make migration name=add_something"; exit 1)
	@n=$$(ls internal/store/migrations | wc -l | tr -d " "); n=$$(printf "%05d" $$((n+1))); \
	 f=internal/store/migrations/$${n}_$(name).sql; \
	 printf -- "-- +goose Up\n\n\n-- +goose Down\n\n" > $$f; echo "created $$f"

.PHONY: dev web up down test migration
