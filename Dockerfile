# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy only the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Set the Go build cache directory
ENV GOCACHE=/root/.cache/go-build

# Enable CGO
ENV CGO_ENABLED=1

# Install GCC (if you need it for any dependencies)
RUN apk add --no-cache build-base

# Enable caching for the Go build process and specify the output binary path
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /main/app .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Copy the executable from the builder stage
COPY --from=builder /main/app /main/app

# Expose a port (you can specify this in your Kubernetes Deployment)
EXPOSE 8000

# Set the command to run the executable
ENTRYPOINT ["/main/app"]
