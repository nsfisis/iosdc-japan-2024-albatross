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
sqldef-dryrun:
	docker compose run --no-TTY tools psqldef --dry-run < ./backend/schema.sql

.PHONY: sqldef
sqldef:
	docker compose run --no-TTY tools psqldef < ./backend/schema.sql
