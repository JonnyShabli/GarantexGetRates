# Используем официальный образ Go как базовый
FROM golang:1.23.2-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /go/src/ggr

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Собираем приложение
RUN go build -o /app ./cmd/main.go

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/ /app

COPY .env /

#EXPOSE 50051
EXPOSE 8080:8080

# Запускаем приложение
CMD ["/app"]
