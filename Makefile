.PHONY: compose-build
compose-build:
	@docker compose build

.PHONY: compose-rm
compose-rm:
	@docker image rm postgres:16.2-alpine
	@docker image rm app:1.0
	@docker volume rm film-library_pgdata

.PHONY: compose-up
compose-up:
	@docker compose up

.PHONY: compose-down
compose-down:
	@docker compose down