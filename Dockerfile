# ----------- Builder Stage -----------
# Use official Go image for building the binary
FROM golang:1.25-alpine AS builder

# Install git (required for go mod) and air for dev
RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Cache go mod deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install swag CLI tool
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Install sql-migrate and build seed-admin for container init
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go build -o /app/bin/seed-admin ./cmd/seed-admin

# Generate Swagger documentation
RUN swag init -g cmd/main.go --output docs --parseDependency


# Build a statically linked binary for production
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/sihp-be -ldflags="-s -w" ./cmd/main.go

# ----------- Development Stage -----------
FROM golang:1.25-alpine AS dev
WORKDIR /app

# Install air, git, and postgres client (pg_isready for entrypoint)
RUN apk add --no-cache git curl postgresql-client
RUN go install github.com/air-verse/air@latest

COPY --from=builder /app /app
COPY --from=builder /go/bin/sql-migrate /usr/local/bin/sql-migrate
COPY --from=builder /app/bin/seed-admin /usr/local/bin/seed-admin
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port for the app
EXPOSE 8080

# Migrate schema, seed admin, then start the app
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]

# ----------- Production Stage -----------
FROM gcr.io/distroless/static-debian12 AS prod

# Create non-root user
USER nonroot:nonroot

WORKDIR /app

# Copy the statically built binary from builder
COPY --from=builder /app/bin/sihp-be /app/sihp-be

# Copy Swagger docs
COPY --from=builder /app/docs /app/docs

# Expose port for the app
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/sihp-be"]
