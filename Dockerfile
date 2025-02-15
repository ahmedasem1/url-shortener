FROM golang:1.22 AS builder
WORKDIR /app

# Copy module files first to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy && go mod download

# Copy the entire source code
COPY . .

# Build the application with Linux target and static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shortener ./src

# ---- Final Minimal Runtime Stage ----
FROM alpine:latest
WORKDIR /app

# Install dependencies (optional, if needed)
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/url-shortener .

# Ensure execute permission
RUN chmod +x url-shortener

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./url-shortener"]
# CMD ["go", "test", "./..."]

