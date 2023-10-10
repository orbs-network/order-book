LOG_LEVEL ?= info

-include .makerc
export

start:
	@echo "Starting the application..."
	@go run cmd/order-book/main.go