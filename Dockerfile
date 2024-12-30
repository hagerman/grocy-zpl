# Build-Stage
FROM golang:1.23-alpine AS build
WORKDIR /app

# Copy the source code
COPY . .

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main

# Deploy-Stage
FROM alpine:3.20.2
WORKDIR /app

# Set environment variable for runtime
ENV GO_ENV=production

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port your application runs on
EXPOSE 8000

# Command to run the application
CMD ["./main"]