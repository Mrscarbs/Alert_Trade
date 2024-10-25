# Build stage
FROM rust:1.75-slim-bullseye as builder

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    pkg-config \
    libssl-dev \
    default-libmysqlclient-dev \
    && rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /usr/src/myapp

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

# Runtime stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    default-mysql-client \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -s /bin/false appuser

# Copy the binary from builder
COPY --from=builder /usr/src/myapp/target/release/MpMlaAPI /usr/local/bin/

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
    CMD curl -f http://localhost:8080/health || exit 1

# Add metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL version="1.0"
LABEL description="MpMla API Server"

# Switch to non-root user
USER appuser

# Create and set working directory
WORKDIR /app

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["MpMlaAPI"]