# Build stage
FROM rust:1.75-slim-bullseye as builder

# Install required dependencies for building
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    pkg-config \
    default-libmysqlclient-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create a new empty shell project
WORKDIR /usr/src/app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Cache dependencies
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    cargo init && \
    cargo build --release && \
    rm src/*.rs

# Copy source code
COPY ./src ./src

# Build for release with dependency caching
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    cargo build --release

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    default-libmysqlclient-dev \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -u 1000 -U app

# Create directory for logs with proper permissions
RUN mkdir /app_logs && chown app:app /app_logs

# Copy the binary from builder
COPY --from=builder /usr/src/app/target/release/api_server /usr/local/bin/additional_rust_api

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
    CMD curl -f http://localhost:8080/health || exit 1

# Add metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL version="1.0"
LABEL description="Additional Rust API Server"

# Switch to non-root user
USER app
WORKDIR /app_logs

# Environment variables
ENV RUST_LOG=info
ENV DATABASE_URL=mysql://root:Karma100%@db/alert_trade_db

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["additional_rust_api"]