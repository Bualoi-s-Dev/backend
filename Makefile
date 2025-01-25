.PHONY: run tidy swag

run: tidy swag
	@echo "Starting the server..."
	go run ./cmd/main.go

tidy:
	@echo "Tidying up Go modules..."
	go mod tidy

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main.go -o ./docs