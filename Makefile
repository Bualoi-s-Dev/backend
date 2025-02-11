.PHONY: run tidy swag server tsgen

run: swag tidy server
	
tidy:
	@echo "Tidying up Go modules..."
	go mod tidy

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main/main.go -o ./docs

server:
	@echo "Starting the server..."
	go run ./cmd/main/main.go

tsgen:
	@echo "Generating TypeScript types..."
	go run ./cmd/tsgen/main.go