# Используем официальный образ Golang в качестве базового
FROM golang:1.22-alpine

# Устанавливаем зависимости для Kafka
RUN apk add --no-cache git bash

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Загружаем модули Go
RUN go mod tidy

# Копируем скрипт wait-for-it.sh в контейнер
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Сборка приложения
RUN go build -o main ./cmd

# Устанавливаем команду запуска контейнера
CMD ["/bin/sh", "-c", "/wait-for-it.sh postgres:5432 -- /wait-for-it.sh kafka1:29092 -- ./main"]