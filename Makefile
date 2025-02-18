# Description: Makefile for the project
# Usage:
# 	make run: Setup Go and run the server
# 	make run-test: Setup Go and run the server in httptest and run tests
# 	make tidy: Tidy up Go modules
# 	make swag: Generate Swagger documentation
# 	make server: Start the server
# 	make testing: Run tests
# 	make tsgen: Generate TypeScript types

.PHONY: run tidy swag server tsgen testing run-test

run: swag tidy server

run-test: swag tidy testing
	
tidy:
	@echo "Tidying up Go modules..."
	go mod tidy

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main/main.go -o ./docs

server:
	@echo "Starting the server..."
	go run ./cmd/main/main.go

testing:
	@echo "Running tests..."
	go test -v ./testing/runner

tsgen:
	@echo "Generating TypeScript types..."
	go run ./cmd/tsgen/main.go