# ===== Build stage =====
FROM golang:1.25.4 AS builder

# Рабочая директория
WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# ===== Run stage =====
FROM alpine:3.19

# Рабочая директория
WORKDIR /app

# Копируем скомпилированное приложение
COPY --from=builder /app/app .

# Копируем шаблоны и статические файлы
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Создаём папку для загруженных файлов
RUN mkdir -p /app/uploads

# Создаём точку монтирования volume для базы данных
RUN mkdir -p /data

# Открываем порт, который использует Go-сервер
EXPOSE 8080

# Команда запуска приложения
CMD ["./app"]
