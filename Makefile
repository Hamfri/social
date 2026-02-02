# make does not support quotes, export ENV_ATTR=123 spaces around `=` in  .env
include .env

MIGRATIONS_PATH=./migrations

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@printf "Are you sure? [y/n] " && read ans && [ "$${ans:-n}" = y ]

.PHONY: migrate/create
migrate/create:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=$(MIGRATIONS_PATH) ${name}

.PHONY: migrate/up
migrate/up: confirm
	@echo 'Running migrations'
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_DSN) up

.PHONY: migrate/down
migrate/down:
	@echo 'Running migrations'
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_DSN) down

