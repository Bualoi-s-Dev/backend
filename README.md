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

This project has used [swaggo](https://github.com/swaggo/swag) to generate API documentation  
The swagger UI run at `localhost:8080/swagger/index.html`

update swagger by

```
swag init -g ./cmd/main.go -o ./docs
```

## Typescript type generation

This project has used [typescriptify](https://github.com/tkrajina/typescriptify-golang-structs) to generate types for TypeScript

generate by

```
go run ./cmd/tsgen/tsgen.go
```

the generated file will be stored in `/gen/api_types.ts`
