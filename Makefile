.PHONY: run tidy swag

run: swag tidy
	@echo "Starting the server..."
	go run ./cmd/main/main.go

tidy:
	@echo "Tidying up Go modules..."
	go mod tidy

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main/main.go -o ./docs