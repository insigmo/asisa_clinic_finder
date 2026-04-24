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

> ℹ️ SQLite база данных хранится в Docker volume `bot_data` и переживает перезапуски контейнера.

---

## 🔧 Управление

| Команда | Описание |
|---------|----------|
| `docker compose up -d` | Запустить в фоне |
| `docker compose down` | Остановить и удалить контейнер |
| `docker compose down -v` | Остановить и **удалить данные** (volume) |
| `docker compose logs -f` | Следить за логами в реальном времени |
| `docker compose restart bot` | Перезапустить бота |
| `docker compose build --no-cache` | Пересобрать образ с нуля |

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

## 🗂 Структура проекта

```
asisa_clinic_finder/
├── cmd/
│   ├── bot/              # Точка входа бота
│   └── data_migrator/    # Загрузчик справочных данных (направления, города)
├── internal/
│   ├── constants/        # Состояния FSM
│   ├── db/               # Работа с SQLite
│   ├── fsmstate/         # Машина состояний (колбэки)
│   ├── handlers/         # Telegram-хендлеры
│   ├── helpers/          # HTTP-клиент ASISA, утилиты
│   ├── i18n/             # Локализация (ru / es / en)
│   ├── keyboards/        # Reply-клавиатуры
│   ├── logger/           # Настройка zap-логгера
│   ├── middleware/        # Telegram middleware (inject DB)
│   ├── model/            # Общие модели данных
│   └── services/clinic/  # Бизнес-логика поиска клиник
├── migrations/           # SQL-миграции (goose)
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── Taskfile.yml
```

---

## ⚙️ Переменные окружения

| Переменная | Обязательная | Описание |
|------------|:------------:|----------|
| `BOT_TOKEN` | ✅ | Токен Telegram-бота от @BotFather |

---

## 🤖 Команды бота

| Команда / Кнопка | Описание |
|-----------------|----------|
| `/start` | Начало работы, ввод города |
| **Найти поликлинику** / *Buscar clínica* / *Find clinic* | Поиск клиники по направлению |
| **Поменять город** / *Cambiar ciudad* / *Change city* | Изменить сохранённый город |
| **Поменять язык** / *Cambiar idioma* / *Change language* | Переключить язык интерфейса |

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

---

## 📦 Docker-образ

- **Builder**: `golang:1.24-bookworm` — полная среда с GCC для сборки CGO (SQLite)
- **Runtime**: `debian:trixie-slim` (Debian 13) — минимальный образ ~80 МБ
- Итоговый образ содержит только бинарь, goose, task и миграции
