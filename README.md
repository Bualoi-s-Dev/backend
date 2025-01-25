# PhotoMatch Backend

## Prerequisites
- Go / Gin

## Building

### Using MakeFile (make required)
```
make run
```

OR

### Manual run

Firstly, create go.tidy

```
go mod tidy
```

Then, create docs for swaggerUI

```
swag init -g ./cmd/main.go -o ./docs
```

Finally, run at localhost:8080

```
go run ./cmd/main.go
```

## Swagger UI

This project has used swagger UI to generate API documentation  
The swagger UI run at `localhost:8080/swagger/index.html`
