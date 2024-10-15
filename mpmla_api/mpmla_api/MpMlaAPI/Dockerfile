# Build stage
FROM rust:latest as builder
WORKDIR /usr/src/myapp
COPY . .

RUN cargo build --release

# Runtime stage
FROM debian:buster-slim

# Install MySQL client, OpenSSL, and CA certificates
RUN apt-get update && \
    apt-get install -y default-mysql-client openssl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /usr/src/myapp/target/release/MpMlaAPI /usr/local/bin/

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["MpMlaAPI"]