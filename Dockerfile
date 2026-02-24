# Build
FROM golang:1.25 as builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# -a ignore cache and rebuild everything
RUN go build -a -installsuffix cgo -o api cmd/api/*.go

# Run
FROM scratch
WORKDIR /app
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /app/api .
ENV APP_PORT=:9000
EXPOSE 9000
CMD ["./api"]