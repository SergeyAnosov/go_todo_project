# Stage 1: Build executable
FROM golang:1.23.4 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-server

# Stage 2: Minimal runtime image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates && \
    adduser -D myuser && \
    mkdir -p /app/db && \
    chown -R myuser /app

WORKDIR /app

# Копируем бинарник и статику
COPY --from=builder /todo-server ./todo-server
COPY web ./web
COPY .env ./.env

# Делаем исполняемый файл запускаемым
RUN chmod +x ./todo-server

# Запуск сервера
CMD ["./todo-server"]