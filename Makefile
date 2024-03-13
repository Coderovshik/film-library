.PHONY: generate-openapi
generate-openapi:
	@oapi-codegen -generate std-http -package openapi api/openapi.yml > pkg/openapi.gen.go

.PHONY: build-migrator
build-migrator:
	@go build -o bin/migrator cmd/migrator/main.go

.PHONY: run-migrator
run-migrator: build-migrator
	@./bin/migrator

PHONY: postgres-init
postgres-init:
	@docker run --name db-test -p 5433:5432 \
		-e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:16-alpine

.PHONY: postgres-rm
postgres-rm:
	@docker stop db-test && docker rm db-test

.PHONY: postgres-cli
postgres-cli:
	@docker exec -it db-test psql -U admin

.PHONY: compose-build
compose-build:
	@docker compose build

.PHONY: compose-rm
compose-rm:
	@

.PHONY: compose-up
compose-up:
	@docker compose up

.PHONY: compose-down
compose-down:
	@docker compose down