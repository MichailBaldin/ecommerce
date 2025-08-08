.PHONY: build up down logs test clean run-user run-product health health-user health-product

# Docker commands
build:
	docker-compose build

up:
	docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v

# Test commands
test:
	cd services/users && go test ./...
	cd services/products && go test ./...

test-user:
	cd services/users && go test ./...

test-product:
	cd services/products && go test ./...

# Local run commands (using separate local databases)
run-user:
	cd services/users && \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_DB=users \
	REDIS_ADDR=localhost:6379 \
	PORT=8001 \
	go run main.go

run-product:
	cd services/products && \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5433 \
	POSTGRES_DB=products \
	REDIS_ADDR=localhost:6380 \
	PORT=8002 \
	go run main.go

# Health check commands
health: health-user health-product

health-user:
	@echo "Checking user service health..."
	@curl -s http://localhost:8001/health || echo "User service not available"

health-product:
	@echo "Checking product service health..."
	@curl -s http://localhost:8002/health || echo "Product service not available"

# API testing commands
test-create-product:
	curl -X POST http://localhost:8002/api/v1/products \
	  -H "Content-Type: application/json" \
	  -d '{"name":"Gaming Laptop","description":"High-performance laptop","price":1299.99}'

test-get-product:
	curl -s http://localhost:8002/api/v1/products/1