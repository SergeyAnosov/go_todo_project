# Stage 1: Build executable
FROM golang:1.23.4 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-server

# Stage 2: Minimal runtime image
FROM ubuntu:latest
WORKDIR /app

# Установка минимальных зависимостей
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Копируем бинарник и статику
COPY --from=builder /todo-server ./todo-server
COPY web ./web
COPY .env ./.env

# Создаем папку для БД
RUN mkdir -p /db

# Делаем исполняемый файл запускаемым
RUN chmod +x ./todo-server

# Передаем переменные окружения через docker-compose
EXPOSE 7540

# Запуск сервера
CMD ["./todo-server"]