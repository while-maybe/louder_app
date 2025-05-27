FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
# -o specifies the output file name
# ./cmd/app/main.go is the path to your main package
# CGO_ENABLED=0 for a statically linked binary (good for small alpine images)
# -ldflags="-w -s" strips debug symbols and DWARF info to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/louder ./cmd/app/main.go
# If you have multiple 'main' packages for different servers (e.g., Gin, Echo):
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/gin-server ./cmd/gin_server/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/echo-server ./cmd/echo_server/main.go

# FINAL STAGE
# Use a minimal base image like alpine for the final image
FROM alpine:latest AS runtime

# Copy the Pre-built binary file from the "builder" stage
COPY --from=builder /app/louder /app/louder
# If you built multiple binaries:
# COPY --from=builder /app/gin-server /app/gin-server
# COPY --from=builder /app/echo-server /app/echo-server

# (Optional) If your application needs CA certificates for HTTPS calls
# RUN apk --no-cache add ca-certificates

# (Optional) Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose port (e.g., 8080 if your Gin/Echo server listens on it)
# This is documentation; you still need to map it with -p when running
EXPOSE 8484

# Command to run the executable
# This will be the default command when the container starts.
# Choose the binary you want to run by default.
ENTRYPOINT ["/app/louder"]
# Or if you want to select which server to run via an environment variable or command override:
# CMD ["/app/my-hexagonal-app"] # or ["/app/gin-server"] or ["/app/echo-server"]