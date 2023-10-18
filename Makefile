LOG_LEVEL ?= info

-include .makerc
export

start:
	@echo "Starting server and db..."
	@docker-compose up --build -d

stop:
	@echo "Stopping server and db..."
	@docker-compose down

restart: stop start

logs: 
	@docker-compose logs -f

test:
	@echo "Running tests..."
	@go test -v ./...