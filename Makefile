.PHONY: ci-local test build clean help workspace-sync

# Локальная эмуляция CI
ci-local:
	@echo "Running local CI checks..."
	@echo "Syncing workspace..."
	go work sync
	@echo "Running go mod tidy..."
	cd services/gateway && go mod tidy
	cd services/users && go mod tidy  
	cd services/products && go mod tidy
	@echo "Running go vet..."
	cd services/gateway && go vet ./...
	cd services/users && go vet ./...
	cd services/products && go vet ./...
	@echo "Checking formatting..."
	@UNFORMATTED=$$(find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -s -l); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@echo "Running tests..."
	cd services/gateway && go test -v ./...
	cd services/users && go test -v ./...
	cd services/products && go test -v ./...
	@echo "Building services..."
	cd services/gateway && go build .
	cd services/users && go build .
	cd services/products && go build .
	@echo "All CI checks passed!"

# Синхронизация workspace
workspace-sync:
	go work sync

# Запуск тестов
test:
	@echo "Running tests..."
	cd services/gateway && go test -v ./...
	cd services/users && go test -v ./...
	cd services/products && go test -v ./...

# Сборка всех сервисов
build:
	@echo "Building all services..."
	mkdir -p bin
	cd services/gateway && go build -o ../../bin/gateway .
	cd services/users && go build -o ../../bin/users .
	cd services/products && go build -o ../../bin/products .
	@echo "Build complete! Binaries in ./bin/"

# Запуск конкретного сервиса локально
run-gateway:
	cd services/gateway && go run .

run-users:
	cd services/users && go run .

run-products:
	cd services/products && go run .

# Форматирование кода
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	cd services/gateway && go mod tidy
	cd services/users && go mod tidy
	cd services/products && go mod tidy

# Линтер
lint:
	golangci-lint run

# Очистка
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	cd services/gateway && go clean ./...
	cd services/users && go clean ./...
	cd services/products && go clean ./...

# Помощь
help:
	@echo "Available commands:"
	@echo "  ci-local    - Run all CI checks locally"
	@echo "  test        - Run all tests"
	@echo "  build       - Build all services"
	@echo "  fmt         - Format code and tidy modules"
	@echo "  lint        - Run golangci-lint"
	@echo "  clean       - Clean build artifacts"
	@echo "  run-*       - Run specific service"