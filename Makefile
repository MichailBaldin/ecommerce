.PHONY: build up down logs test clean

build:
	docker-compose build

up:
	docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f

test:
	cd services/users && go test ./...

clean:
	docker-compose down -v

run-user:
	cd services/users && \
	POSTGRES_HOST=localhost \
	REDIS_ADDR=localhost:6379 \
	PORT=8001 \
	go run main.go

health:
	curl -s http://localhost:8001/health