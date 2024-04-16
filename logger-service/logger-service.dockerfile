FROM alpine:latest

# Create app directory
RUN mkdir /app

# Copy binary from builder stage
COPY loggerServiceApp /app

# Set working directory
WORKDIR /app

# Command to run the binary
CMD ["./loggerServiceApp"]

# Add metadata as a label
LABEL authors="rossado"
