# Subscription Service

REST API сервис для управления подписками пользователей.

## Описание

Сервис предоставляет CRUDL-операции для подписок, а также возможность подсчёта общей стоимости подписок по фильтрам.

## Технологии

- Go 1.XX
- Gin Web Framework
- PostgreSQL
- Swagger (swaggo) для автогенерации документации
- Logrus для логирования
- Docker Compose для сборки и запуска

## Установка и запуск

1. Клонировать репозиторий

```bash
git clone <URL>
cd subscription-service
```

2. Создать и настроить базу данных PostgreSQL.

3. Запустить миграции (если есть).

4. Запустить сервис через Docker Compose:

```bash
docker-compose up --build
```

Или локально:

```bash
go run cmd/app/main.go
```

## Конфигурация

Используется файл `.env` или `.yaml` с настройками подключения к базе и портом.

## API документация

Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

## Основные эндпоинты

| Метод | Путь                     | Описание                       |
|-------|--------------------------|--------------------------------|
| POST  | /subscriptions           | Создать подписку               |
| GET   | /subscriptions/{id}      | Получить подписку по ID        |
| GET   | /subscriptions           | Получить все подписки          |
| PUT   | /subscriptions/{id}      | Обновить подписку              |
| DELETE| /subscriptions/{id}      | Удалить подписку               |
| GET   | /subscriptions/total     | Подсчитать суммарную стоимость |

## Пример запроса создания подписки

```bash
curl -X POST http://localhost:8080/subscriptions   -H "Content-Type: application/json"   -d '{
    "service_name": "Netflix",
    "price": 10,
    "user_id": "user-uuid",
    "start_date": "2025-07-01",
    "end_date": "2025-07-31"
}'
```

## Логирование

Используется logrus, логи выводятся в stdout.

## Swagger генерация

Для генерации Swagger документации используется swag:

```bash
swag init -g cmd/app/main.go -o docs
```

## Лицензия

MIT