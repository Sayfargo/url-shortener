include .env
export

service-up:
	docker compose up -d

service-down:
	docker compose down -v

goose-create:
	@if [ -z "$(name)" ]; then \
		echo "There is no name param. Example: make migrate-create name=name_of_migrate"; \
		exit 1; \
	fi;	\

	docker compose run --rm \
	-e GOOSE_COMMAND="create" \
	-e GOOSE_COMMAND_ARG="$(name) sql" \
	migrations

goose-status:
	docker compose run --rm \
	-e GOOSE_COMMAND="status" \
	migrations

goose-up:
	docker compose run --rm \
	-e GOOSE_COMMAND="up" \
	migrations

goose-up-by-one:
	docker compose run --rm \
	-e GOOSE_COMMAND="up-by-one" \
	migrations

goose-down:
	docker compose run --rm \
	-e GOOSE_COMMAND="down" \
	migrations

goose-down-to:
	@if [ -z "$(id)" ]; then \
		echo "There is no migrate id. Example: make goose-down-to id=20170614145246"; \
		exit 1; \
	fi
	docker compose run --rm \
		-e GOOSE_COMMAND="down-to"
		-e GOOSE_COMMAND_ARG="$(id)" \
		migrations

goose-up-to:
	@if [ -z "$(id)" ]; then \
		echo "There is no migrate id. Example: make goose-down-to id=20170614145246"; \
		exit 1; \
	fi
	docker compose run --rm \
		-e GOOSE_COMMAND="up-to"
		-e GOOSE_COMMAND_ARG="$(id)" \
		migrations

goose-full-rollback:
	docker compose run --rm migrations down-to 0 