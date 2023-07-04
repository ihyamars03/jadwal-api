# Stage 1: Build
FROM golang:1.17-alpine AS build

# Set working directory
WORKDIR /app


# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o app

# Stage 2: Final Image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/app .


# Expose the port (if needed)
EXPOSE 3030

# Run the binary
CMD ["./app"]
