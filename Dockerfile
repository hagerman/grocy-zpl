# Build-Stage
FROM golang:1.23-alpine AS builder

# Set Environment Variables
ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build app
RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy any .zpl files from the project directory to the output directory
COPY --from=builder /app/*.zpl .

EXPOSE 8000

CMD [ "./main" ]