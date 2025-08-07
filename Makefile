.PHONY: ci-local test build clean

# Локальная эмуляция CI
ci-local:
	@echo "Running local CI checks..."
	go work sync
	find . -name go.mod -execdir go mod tidy \;
	go vet ./...
	gofmt -s -l . | tee /tmp/gofmt.out && test ! -s /tmp/gofmt.out
	go test ./...
	@echo "Building all services..."
	cd services/user-service && go build .
	cd services/product-service && go build .
	cd services/api-gateway && go build .
	@echo "All CI checks passed!"

# Запуск тестов
test:
	go work sync
	go test -v ./...

# Сборка всех сервисов
build:
	cd services/user-service && go build -o bin/user-service .
	cd services/product-service && go build -o bin/product-service .
	cd services/api-gateway && go build -o bin/api-gateway .

# Очистка
clean:
	find . -name "bin" -type d -exec rm -rf {} +
	go clean ./...