FROM golang:1.24-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка
RUN go build -o main ./cmd/app/main.go

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]
