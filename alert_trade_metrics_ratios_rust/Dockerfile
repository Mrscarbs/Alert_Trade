# Use the latest official Rust image as a parent image
FROM rust:latest as builder

# Set the working directory in the container
WORKDIR /usr/src/myapp

# Copy the entire project
COPY . .

# Build the application
RUN cargo build --release

# Use a newer base image that includes libssl.so.3
FROM debian:bookworm-slim

# Install OpenSSL and CA certificates
RUN apt-get update && apt-get install -y libssl3 ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /usr/src/myapp/target/release/alert_trade_metrics_ratios_rust /usr/local/bin/

# Expose port 8000
EXPOSE 8000

# Set the startup command
CMD ["alert_trade_metrics_ratios_rust"]