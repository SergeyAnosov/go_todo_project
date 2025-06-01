FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY todo-server ./todo-server
COPY web ./web

RUN chmod +x ./todo-server

EXPOSE 8080

# Команда запуска приложения
CMD ["./todo-server"]