FROM golang:latest AS builder

COPY . /app
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go mod tidy
RUN swag init -g ./cmd/main.go -o ./docs
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin /app/cmd/main.go

FROM alpine:latest

COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

EXPOSE 8080

CMD ["./appbin"]