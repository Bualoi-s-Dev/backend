# PhotoMatch Backend

## Prerequisites
- Go / Gin

## Building

### Using MakeFile (make required)
```
make run
```

OR

### Docker-compose (docker required)

```
docker-compose up --build -d
```

## Swagger UI

This project has used swagger UI to generate API documentation  
The swagger UI run at `localhost:8080/swagger/index.html`

update swagger by

```
swag init -g ./cmd/main.go -o ./docs
```
