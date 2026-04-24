FROM golang:1.26.2-bookworm AS builder

WORKDIR /app

COPY ./go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/bot ./cmd/bot

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/migrator ./cmd/data_migrator


FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /app/db /app/bin

COPY ./ /app/
COPY --from=builder /app/bin/bot           /app/bot
COPY --from=builder /app/bin/migrator      /app/migrator
COPY --from=builder /go/bin/goose          /usr/local/bin/goose
COPY --from=builder /app/migrations        /app/migrations
