.PHONY: build
build:
	docker compose build

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: psql
psql: up
	docker compose exec db psql --user=postgres albatross

.PHONY: sqldef-dryrun
sqldef-dryrun: down build
	docker compose up -d db
	docker compose run --no-TTY tools psqldef --dry-run < ./backend/schema.sql

.PHONY: sqldef
sqldef: down build
	docker compose up -d db
	docker compose run --no-TTY tools psqldef < ./backend/schema.sql

.PHONY: oapi-codegen
oapi-codegen:
	cd backend; make oapi-codegen

.PHONY: sqlc
sqlc:
	cd backend; make sqlc
