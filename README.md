# go-subscription-service

![Status](https://img.shields.io/badge/status-alpha-blue)
![Go](https://img.shields.io/badge/go-%3E%3D1.20-brightgreen)
![Docker](https://img.shields.io/badge/docker-ready-blue)
![Postgres](https://img.shields.io/badge/postgres-required-blue)

## Кратко

`go-subscription-service` — REST API на Go для управления подписками.
Проект включает:

* REST API (JSON),
* Swagger/OpenAPI документацию,
* хранение данных в PostgreSQL,
* структурированное и централизованное логирование,
* контейнеризацию с помощью Docker (готовые Dockerfile / docker-compose).

---

# Возможности

* CRUDL для сущности `Subscription`
* Swagger UI для интерактивной документации
* Подключение к PostgreSQL (через пул соединений)
* Централизованное логирование в файл(структурированные логи)
* Конфигурация через переменные окружения
* Собирается и запускается в Docker-контейнере

---

# Быстрый старт (Docker)

1. Создайте файл `.env` (пример ниже).
2. Запустите PostgreSQL и сервис в Docker:

```bash
# Запустить через docker-compose (пример)
docker-compose up -d --build
```

3. Проверьте логи сервиса:

```bash
docker-compose logs -f service
```

4. Откройте Swagger UI:

```
http://localhost:8080/swagger/index.html
```

> Порт `8080` и имена контейнеров — примеры. Замените на ваши значения.

---

# Пример `docker-compose.yml`

```yaml
services:
  app:
    build: .
    container_name: subscription_service
    env_file:
      - .env
    ports:
      - "${APP_PORT}:${APP_PORT}"
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    container_name: subscription__service_db
    env_file:
      - .env
    ports:
      - "${DB_PORT}:5432"
    volumes:
      - db_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  db_data:
```

---

# Пример `.env`

```env
APP_PORT=1111

# settings for service
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASS=1111
DB_NAME=user

#container settings
POSTGRES_DB=database
POSTGRES_USER=user
POSTGRES_PASSWORD=1111
```

---

# Локальная сборка и запуск (Go)

```bash
# Установите зависимости
go mod download

# Сборка бинарника
go build -o bin/subscription-service ./cmd/server

# Запуск (через .env или переменные окружения)
./bin/subscription-service
```

---

# API (пример)

> Ниже — примерная схема endpoint-ов. Подставьте реальные пути и схемы из Swagger.

* `GET /api/v1/subscriptions` — список подписок
* `GET /api/v1/subscriptions/{id}` — получить подписку
* `POST /api/v1/subscriptions` — создать подписку
* `PUT /api/v1/subscriptions/{id}` — обновить подписку
* `DELETE /api/v1/subscriptions/{id}` — удалить подписку

Пример запроса: создание подписки

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
  "service_name": "Sberbank",
  "price": 1099.0,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "2025-11-10T00:00:00Z",
  "end_date": "2026-11-10T00:00:00Z"
}'
```