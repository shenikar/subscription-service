# 1. Базовый образ Go
FROM golang:1.24-alpine

# 2. Установка зависимостей
RUN apk add --no-cache git

# 3. Рабочая директория внутри контейнера
WORKDIR /app

# 4. Копируем go.mod и go.sum для кеширования
COPY go.mod go.sum ./
RUN go mod download

# 5. Копируем остальной код
COPY . .

# 6. Сборка
RUN go build -o main ./cmd/main.go

# 7. Установка переменной окружения
ENV GIN_MODE=release

# 8. Порт приложения
EXPOSE 8080

# 9. Команда по умолчанию
CMD ["./main"]
