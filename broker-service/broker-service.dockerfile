FROM alpine:latest

# Create app directory
RUN mkdir /app

# Copy binary from builder stage
COPY brokerApp /app

# Set working directory
WORKDIR /app

# Command to run the binary
CMD ["./brokerApp"]

# Add metadata as a label
LABEL authors="rossado"
