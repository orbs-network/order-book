LOG_LEVEL ?= info

-include .makerc
export

start:
	@echo "Starting server and db..."
	@docker-compose up --build -d

stop:
	@echo "Stopping server and db..."
	@docker-compose down

watch: start
	@echo "Watching for file changes..."
	@docker-compose watch

restart: stop start

logs: 
	@docker-compose logs -f

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	golangci-lint run

coverage:
	@echo "Generating test coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out