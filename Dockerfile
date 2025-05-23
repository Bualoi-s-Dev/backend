# Build stage
FROM golang:alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git curl

# Set up working directory
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Install Swaggo and generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main/main.go -o ./docs

# Install dependencies & build the binary
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin /app/cmd/main/main.go

# Download Swagger UI
RUN mkdir -p /app/swagger-ui && \
    wget -qO- https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.11.0.tar.gz | tar xz --strip-components=1 -C /app/swagger-ui

# Production stage
FROM alpine:latest

RUN apk add --no-cache tzdata
# Create a non-root user
RUN adduser -D -g '' appuser

# Copy the Go binary and Swagger UI files
COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

# Set permissions
RUN chown -R appuser:appuser /home/appuser/app
USER appuser

RUN ls -la /

EXPOSE 8080

# Run the Go server directly
CMD ["/home/appuser/app/appbin"]