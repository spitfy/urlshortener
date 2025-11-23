# Makefile для URL Shortener

.PHONY: test test-coverage test-coverage-html build clean

# Переменные
COVERAGE_FILE = coverage.out
BINARY_NAME = urlshortener

# Запуск всех тестов
test:
	go test -v ./...

# Запуск тестов с проверкой покрытия
test-coverage:
	go test -coverprofile=$(COVERAGE_FILE) ./internal/...
	go tool cover -func=$(COVERAGE_FILE)

# Детальный HTML отчет о покрытии
test-coverage-html:
	go test -coverprofile=$(COVERAGE_FILE) ./internal/...
	go tool cover -html=$(COVERAGE_FILE)

# Показать только общий процент покрытия
test-coverage-total:
	@go test -coverprofile=$(COVERAGE_FILE) ./... > /dev/null 2>&1
	@go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}'

# Проверка что покрытие выше 70%
test-coverage-check:
	@coverage=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
    result=$$(echo "$$coverage < 70" | bc); \
    if [ $$result -eq 1 ]; then \
        echo "❌ Coverage is too low: $$coverage% (min 70%)"; \
        exit 1; \
    else \
        echo "✅ Coverage: $$coverage%"; \
    fi

test-handler:
	go test -v ./internal/handler/...

test-service:
	go test -v ./internal/service/...

# Сборка проекта
build:
	go build -o $(BINARY_NAME) ./cmd/urlshortener

# Очистка
clean:
	rm -f $(BINARY_NAME) $(COVERAGE_FILE)
	rm -f *.out

# Помощь
help:
	@echo "Available targets:"
	@echo "  test			   - Run all tests"
	@echo "  test-coverage	  - Run tests with coverage report"
	@echo "  test-coverage-html - Open HTML coverage report"
	@echo "  test-coverage-check - Check if coverage >= 70%"
	@echo "  test-handler	   - Run handler tests"
	@echo "  test-service	   - Run service tests"
	@echo "  build			  - Build binary"
	@echo "  clean			  - Clean generated files"
	@echo "  help			   - Show this help"