# Start from the official Go 1.22 image
FROM golang:1.23-alpine

# Install MySQL client
RUN apk add --no-cache mysql-client

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the current directory to the working directory in the container
COPY . .

# Download all dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose port 8080 for the API
EXPOSE 8083

# Run the application
CMD ["./main"]