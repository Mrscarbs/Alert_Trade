# Build stage
FROM rust:1.75-slim-bullseye as builder

# Install build dependencies
RUN apt-get update && \
    apt-get install -y default-libmysqlclient-dev libssl-dev pkg-config && \
    rm -rf /var/lib/apt/lists/*

# Copy source code
WORKDIR /app
COPY . .

# Build the application
RUN cargo build --release

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y default-libmysqlclient-dev libssl1.1 ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/target/release/additional_rust_api /usr/local/bin/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["additional_rust_api"]