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

.PHONY: compose-build-up
compose-build-up:
	@docker compose up --build

.PHONY: compose-down
compose-down:
	@docker compose down

.PHONY: test
test:
	@go test -v ./...

.PHONY: test-coverage
test-coverage:
	@go test -cover ./...

.PHONY: test-docker
test-docker:
	@docker build -f Dockerfile -t docker-filmlib-test --progress plain --no-cache --target tester .
	@docker image rm docker-filmlib-test

.PHONY: gen-mock
gen-mock:
	@mockgen -self_package=github.com/Coderovshik/film-library/internal/user -package=user \
		-source=internal/user/user.go -destination=internal/user/mock.go
	@mockgen -self_package=github.com/Coderovshik/film-library/internal/actor -package=actor \
		-source=internal/actor/actor.go -destination=internal/actor/mock.go
	@mockgen -self_package=github.com/Coderovshik/film-library/internal/film -package=film \
		-source=internal/film/film.go -destination=internal/film/mock.go

.PHONY: rm-mock
rm-mock:
	@rm -rf internal/user/mock.go
	@rm -rf internal/actor/mock.go
	@rm -rf internal/film/mock.go