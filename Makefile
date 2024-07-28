.PHONY: build
build:
	docker compose build

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: api-server-only-build
api-server-only-build:
	docker compose build api-server-only

.PHONY: api-server-only-up
api-server-only-up:
	docker compose up -d api-server-only

.PHONY: api-server-only-down
api-server-only-down:
	docker compose down api-server-only db

.PHONY: psql
psql:
	docker compose exec db psql --user=postgres albatross

.PHONY: sqldef-dryrun
sqldef-dryrun: down
	docker compose build db
	docker compose up -d db
	docker compose run --no-TTY tools psqldef --dry-run < ./backend/schema.sql

.PHONY: sqldef
sqldef: down
	docker compose build db
	docker compose up -d db
	docker compose run --no-TTY tools psqldef < ./backend/schema.sql

.PHONY: oapi-codegen
oapi-codegen:
	cd backend; make oapi-codegen

.PHONY: openapi-typescript
openapi-typescript:
	cd frontend; make openapi-typescript

.PHONY: sqlc
sqlc:
	cd backend; make sqlc
