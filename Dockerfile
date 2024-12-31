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
RUN go build ./cmd/grocy-zpl

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/grocy-zpl .

# Copy the contents of /app/assets to /root, preserving the directory structure
COPY /app/assets/ /root/

EXPOSE 8000

CMD [ "./grocy-zpl" ]