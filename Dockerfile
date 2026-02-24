# Build
FROM golang:1.25 AS builder
WORKDIR /app

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# so that swag is available
ENV PATH=$PATH:/go/bin
RUN swag init -g ./api/main.go -d cmd,internal

# -a ignore cache and rebuild everything
RUN go build -a -installsuffix cgo -o api cmd/api/*.go

# Run
# FROM scratch
FROM alpine:latest
WORKDIR /app
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations /app/migrations
ENV APP_PORT=:9000
EXPOSE 9000
CMD ["sh", "-c", "./migrate -path ./migrations -database \"$DB_DSN\" up && ./api"]