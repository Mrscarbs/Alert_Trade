# Use the latest official Rust image as a parent image
FROM rust:latest as builder

# Set the working directory in the container
WORKDIR /usr/src/myapp

# Copy the entire project
COPY . .

# Build the application
RUN cargo build --release

# Use a smaller base image for the final image
FROM debian:buster-slim

# Install OpenSSL and CA certificates
RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*

# Set environment variable for the database URL
# Replace the values with your actual database credentials and connection details
ENV DATABASE_URL=mysql://user:password@host:port/db_name

# Copy the binary from the builder stage
COPY --from=builder /usr/src/myapp/target/release/MpMlaRust /usr/local/bin/

# Set the startup command
CMD ["MpMlaRust"]
