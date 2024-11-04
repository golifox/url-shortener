# Use a specific version of the golang base image
FROM golang:1.23.2 AS builder

# Install necessary packages which includes git and gcc. No need to install bash as the default shell is already bash.
RUN apt-get update && apt-get install -y git gcc

# 'WORKDIR' should have a path where the application source code will reside inside the container
WORKDIR /app

# Copy 'go.mod' and 'go.sum' files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download all dependencies. If the go.mod and go.sum files have not changed, dependencies will be cached.
RUN go mod download

# Copy the source code from the current directory to the '/app' directory inside the container
COPY . .

# Build the Go app. The '-o main' sets the output binary name to 'main'.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Execution stage, start fresh with a smaller base image to reduce final image size
FROM alpine:latest

# Install certificates for SSL/TLS, if you are making secure outbound requests from your application
RUN apk --no-cache add ca-certificates

# Copy the compiled binary 'main' from the 'builder' stage to the new stage
COPY --from=builder /app/main .

# Set environment variables
ENV LINK_LENGTH=7 \
    PUBLIC_URL=http://localhost:3000 \
    REDIS_URL=redis:6379 \
    PORT=3000

# Expose the port the application runs on
EXPOSE $PORT

# Run the binary
CMD ["./main"]