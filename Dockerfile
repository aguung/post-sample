# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies for Atlas
RUN apk add --no-cache curl

# Install Atlas
RUN curl -sSf https://atlasgo.sh | sh

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates libc6-compat tzdata

# Copy from builder
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas
COPY --from=builder /app/migrations ./migrations
COPY entrypoint.sh .

# Make entrypoint executable
RUN chmod +x entrypoint.sh

# Expose port
EXPOSE 8080

# Run entrypoint
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./main", "api"]
