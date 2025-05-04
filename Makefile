.PHONY: build run lint test migrate-up migrate-down docker-up docker-down swagger

# Билд бинарника
build:
	go build -o bin/app ./cmd/main.go

# Запуск приложения
run:
	go run ./cmd/main.go

# Линтер
lint:
	golangci-lint run ./...

# Юнит-тесты
test:
	go test -v ./...

# Миграции вверх (создание таблиц)
migrate-up:
	migrate -path migrations -database $$DATABASE_URL up

# Миграции вниз (откат)
migrate-down:
	migrate -path migrations -database $$DATABASE_URL down

# Запуск docker-compose
docker-up:
	docker-compose up -d

# Остановка docker-compose
docker-down:
	docker-compose down

# Генерация Swagger-документации
swagger:
	swag init -g cmd/main.go -o ./docs
