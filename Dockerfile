FROM golang:latest AS builder

COPY . /app
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main/main.go -o ./docs
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin /app/cmd/main/main.go

# Download Swagger UI
RUN mkdir -p /app/swagger-ui && \
    wget -qO- https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.11.0.tar.gz | tar xz --strip-components=1 -C /app/swagger-ui

FROM alpine:latest

# Install Nginx
RUN apk update && apk add nginx

# Copy the Go binary and Swagger UI files
COPY --from=builder /app /home/appuser/app

# Configure Nginx
COPY nginx.conf /etc/nginx/nginx.conf

WORKDIR /home/appuser/app

EXPOSE 8080

# Start both Nginx and the Go app
CMD ["/bin/sh", "-c", "nginx -g 'daemon off;' & ./appbin"]

