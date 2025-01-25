.PHONY: run tidy

run: tidy
	@echo "Starting the server..."
	go run ./cmd/main.go

tidy:
	@echo "Tidying up Go modules..."
	go mod tidy