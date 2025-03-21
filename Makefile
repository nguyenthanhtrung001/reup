# include ./.env
export
BINARY=engine
## Run the application
run-api:
	@echo "Running the application"
	@go run cmd/api/main.go
## Run the scheduler
run-scheduler:
	@echo "Running the scheduler"
	@go run cmd/scheduler/main.go

## Run the consumer
run-consumer:
	@echo "Running the consumer"
	@go run cmd/consumer/main.go