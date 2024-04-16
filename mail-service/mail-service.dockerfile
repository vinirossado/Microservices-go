FROM alpine:latest

# Create app directory
RUN mkdir /app
# Copy binary from builder stage
COPY mailerApp /app
COPY templates /templates

# Set working directory

# Command to run the binary
CMD ["./app/mailerApp"]

# Add metadata as a label
LABEL authors="rossado"
