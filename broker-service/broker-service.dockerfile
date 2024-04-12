## Use the official Golang image from Docker Hub
#FROM golang:1.22-alpine as builder
#
## Set the working directory inside the container
#WORKDIR /app
#
## Copy the Go module files
#COPY go.mod go.sum ./
#
## Download and install any required dependencies
#RUN go mod download
#
## Copy the entire project directory into the container
#COPY . .
#
## Build the Go app
#RUN go build -o main ./cmd/api
#
## Expose port 8080 to the outside world
#EXPOSE 8080
#
## Command to run the executable
#CMD ["./main"]
#

# Builder stage
FROM golang:1.22-alpine as builder

# Set the GOPROXY environment variable to bypass the need for modules to be fetched directly from VCS
ENV GOPROXY=https://proxy.golang.org

# Set the working directory
WORKDIR /app

COPY go.mod go.sum ./
# Copy the entire project to the working directory
COPY . .

# Initialize Go modules and download dependencies
RUN go mod download

# Build the Go binary
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# Final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/brokerApp .

# Make the binary executable
RUN chmod +x brokerApp

# Define the command to run the binary
CMD ["./brokerApp"]

# Add metadata as a label
LABEL authors="rossado"


# Final stage
FROM alpine:latest

# Create app directory
RUN mkdir /app

# Copy binary from builder stage
COPY --from=builder /app/brokerApp /app/

# Set working directory
WORKDIR /app

# Command to run the binary
CMD ["./brokerApp"]

# Add metadata as a label
LABEL authors="rossado"
