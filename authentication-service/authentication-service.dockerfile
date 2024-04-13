FROM alpine:latest

# Create app directory
RUN mkdir /app

# Copy binary from builder stage
COPY authApp /app/

# Set working directory
WORKDIR /app

# Command to run the binary
CMD ["./authApp"]

# Add metadata as a label
LABEL authors="rossado"
