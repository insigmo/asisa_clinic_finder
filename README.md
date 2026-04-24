# 🏥 ASISA Clinic Finder Bot

Telegram-бот для поиска клиник в сети **ASISA** (Испания) по городу и медицинскому направлению.  
Поддерживает три языка: 🇷🇺 русский, 🇪🇸 испанский, 🇬🇧 английский.

---

## 🚀 Быстрый старт

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/insigmo/asisa_clinic_finder.git
cd asisa_clinic_finder
```

### 2. Создайте файл окружения

```bash
cp .env.example .env
```

Откройте `.env` и вставьте токен вашего бота:

```env
BOT_TOKEN=1234567890:AABBCCDDEEFFaabbccddeeff1234567890
```

### 3. Запустите

```bash
docker compose up -d
```

При первом запуске Docker автоматически:
1. **Собирает образ** из `Dockerfile` на базе `debian:trixie-slim`
2. **Выполняет миграции** (`task migrate`) — создаёт таблицы и загружает справочные данные
3. **Запускает бота**

---

## 🔧 Управление

| Команда                           | Описание                                 |
|-----------------------------------|------------------------------------------|
| `docker compose up -d`            | Запустить в фоне                         |
| `docker compose down`             | Остановить и удалить контейнер           |
| `docker compose down -v`          | Остановить и **удалить данные** (volume) |
| `docker compose logs -f`          | Следить за логами в реальном времени     |
| `docker compose restart bot`      | Перезапустить бота                       |
| `docker compose build --no-cache` | Пересобрать образ с нуля                 |

### Просмотр логов

```bash
docker compose logs -f bot
```

### Выполнить миграции вручную

Если нужно запустить миграции отдельно (например, после добавления новых):

```bash
docker compose run --rm bot task migrate
```

---

## 🏗 Сборка без Docker (локально)

```bash
# Установите зависимости
go mod download

# Выполните миграции (требуется goose и task)
task migrate

# Запустите бота
BOT_TOKEN=<ваш_токен> go run ./cmd/bot
```