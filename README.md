# PhotoMatch Backend

## Prerequisites

Ensure you have the following installed:
- [Go](https://go.dev/) (with [Gin](https://gin-gonic.com/))
- [Make](https://www.gnu.org/software/make/) (optional, for easier commands)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

## Getting Started

### Run the Project

#### Using Makefile (Requires `make`)
```sh
make run
```

#### Using Docker-Compose (Requires `docker`)
```sh
docker-compose up --build -d
```

## API Documentation

This project uses [Swaggo](https://github.com/swaggo/swag) for generating API documentation.

### Access Swagger UI
Swagger UI is available at:
```
http://localhost:8080/swagger/index.html
```

### Update Swagger Documentation
Use one of the following commands to regenerate Swagger documentation:
```sh
make swag
```
Or manually:
```sh
swag init -g ./cmd/main.go -o ./docs
```

## TypeScript Type Generation

The project leverages [typescriptify](https://github.com/tkrajina/typescriptify-golang-structs) to generate TypeScript types from Go structs.

### Generate TypeScript Types
```sh
make tsgen
```
Or manually:
```sh
go run ./cmd/tsgen/main.go
```

The generated TypeScript file will be located at:
```
/gen/api_types.ts
```

## License

This project is licensed under the MIT License. See `LICENSE` for more details.