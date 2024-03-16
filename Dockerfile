FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x
COPY . .

RUN go build -o ./bin/app cmd/filmlib/main.go
RUN go build -o ./bin/migrator cmd/migrator/main.go

FROM builder AS tester

RUN go test -v ./...
RUN go test -cover ./...

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/bin/migrator /

EXPOSE 8080

CMD ["ash", "-c", "/migrator;/app"]