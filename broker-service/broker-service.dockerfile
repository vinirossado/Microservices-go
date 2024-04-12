# Builder stage
FROM golang:1.22-alpine as builder

RUN mkdir /app

# Set the GOPROXY environment variable to bypass the need for modules to be fetched directly from VCS
ENV GOPROXY=https://proxy.golang.org

# Copy only the go.mod and go.sum files to leverage Docker layer caching
COPY go.mod go.sum /app/

WORKDIR /app

# Initialize Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . /app

WORKDIR /app

# Build the Go binary
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# Final stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]

LABEL authors="rossado"
