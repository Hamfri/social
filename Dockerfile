#--------------------------------------------
# Stage 1: Install Tools
#--------------------------------------------
FROM golang:1.25-alpine AS tools

# install swag for generating api docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

#--------------------------------------------
# Stage 2: Builder
#--------------------------------------------
FROM golang:1.25 AS builder
WORKDIR /app

COPY --from=tools /go/bin/swag /go/bin/swag

# Cache Go modules (runs only if go.sum and go.mod changes)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and generate swagger docs
COPY . .
RUN /go/bin/swag init -g ./api/main.go -d cmd,internal

ARG VERSION=dev

# Build the binary
# -s -w reduces binary size for faster container cold starts
# -s strip symbol table. In case of a crash no stack-trace
# -w strip DWARF. Removes debugging information, gdb | delve debbugers can't be used.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X 'main.version=${VERSION}'" -o api cmd/api/*.go

#--------------------------------------------
# Stage 3: Final Runtime (Production)
#--------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN adduser -D social
USER social

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]