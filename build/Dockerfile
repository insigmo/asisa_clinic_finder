# syntax=docker/dockerfile:1
# ┌─────────────────────────────────────────────────────────────────────────────
# │ Stage 1: build
# │ golang:1.24-bookworm содержит gcc/glibc, нужные для CGO (go-sqlite3)
# └─────────────────────────────────────────────────────────────────────────────
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Копируем манифесты зависимостей отдельно — кэшируем слой go mod download
COPY ../go.mod go.sum ./
RUN go mod download

# Устанавливаем Task и goose для выполнения миграций
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем исходный код и собираем бинарь
COPY .. .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/bot ./cmd/bot

# ┌─────────────────────────────────────────────────────────────────────────────
# │ Stage 2: runtime
# │ debian:trixie-slim — последний стабильный Debian 13, минимальный образ
# └─────────────────────────────────────────────────────────────────────────────
FROM debian:trixie-slim AS runtime

# Создаём непривилегированного пользователя
RUN groupadd --gid 1001 appgroup && \
    useradd  --uid 1001 --gid appgroup --no-create-home appuser

WORKDIR /app

# Копируем бинарь бота
COPY --from=builder /app/bin/bot           ./bot
# Копируем task и goose для выполнения миграций внутри контейнера
COPY --from=builder /go/bin/task           /usr/local/bin/task
COPY --from=builder /go/bin/goose          /usr/local/bin/goose
# Копируем миграции и Taskfile
COPY --from=builder /app/migrations        ./migrations
COPY --from=builder /app/Taskfile.yml      ./Taskfile.yml

# Создаём директорию для SQLite базы данных
RUN mkdir -p /app/db && chown -R appuser:appgroup /app

USER appuser

# Переменная окружения — токен бота передаётся снаружи
ENV BOT_TOKEN=""

CMD ["./bot"]
