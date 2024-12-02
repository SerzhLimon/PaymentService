all: run

run:
	@echo "Running Docker containers..."
	docker-compose up -d
	go run cmd/main.go

stop:
	@echo "Stopping Docker containers..."
	docker-compose down

test: 
	go test ./tests -v