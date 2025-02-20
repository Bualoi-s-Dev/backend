# Build stage
FROM golang:latest AS builder

COPY . /app
WORKDIR /app

# Install Swaggo and generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main/main.go -o ./docs

# Install dependencies & build the binary
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin /app/cmd/main/main.go

# Download Swagger UI
RUN mkdir -p /app/swagger-ui && \
    wget -qO- https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.11.0.tar.gz | tar xz --strip-components=1 -C /app/swagger-ui

# Production stage
FROM alpine:latest

# Install Nginx & Supervisor
RUN apk update && apk add nginx supervisor

# Copy the Go binary and Swagger UI files
COPY --from=builder /app /home/appuser/app

# Configure Nginx
COPY nginx.conf /etc/nginx/nginx.conf

# Copy Supervisor config to manage both processes
COPY supervisord.conf /etc/supervisord.conf

WORKDIR /home/appuser/app

EXPOSE 8080

# Run both Nginx and Go Gin using Supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
