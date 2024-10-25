# Builder stage
FROM rust:1.75-slim as builder

# Install required dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    default-libmysqlclient-dev \
    openssl \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create a new empty shell project
WORKDIR /usr/src/app
RUN cargo new --bin api_server
WORKDIR /usr/src/app/api_server

# Copy over your manifests
COPY ./Cargo.toml ./Cargo.toml
COPY ./Cargo.lock ./Cargo.lock

# Cache dependencies
RUN cargo build --release
RUN rm src/*.rs

# Copy your source code
COPY ./src ./src

# Build for release
RUN rm ./target/release/deps/api_server*
RUN cargo build --release

# Runtime stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    default-libmysqlclient-dev \
    openssl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -u 1000 -U app

# Copy the build artifact from builder
COPY --from=builder /usr/src/app/api_server/target/release/api_server /usr/local/bin/

# Create directory for logs with proper permissions
RUN mkdir /app_logs && chown app:app /app_logs

# Switch to non-root user
USER app

# Set the working directory
WORKDIR /app_logs

# Set environment variables
ENV RUST_LOG=info
ENV DATABASE_URL=mysql://root:Karma100%@db/alert_trade_db

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["additional_rust_api"]