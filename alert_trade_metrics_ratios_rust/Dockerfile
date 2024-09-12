# Build stage
FROM rust:1.73 as builder

WORKDIR /usr/src/main

# Copy over manifests
COPY Cargo.toml Cargo.lock ./

# Create a dummy main.rs to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build and cache dependencies
RUN cargo build --release

# Remove the dummy main
RUN rm src/main.rs

# Now copy your actual source code
COPY src ./src

# Build the application
RUN touch src/main.rs && cargo build --release

# Final stage
FROM debian:bullseye-slim

# Install OpenSSL and CA certificates
RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the build artifact from the build stage
COPY --from=builder /usr/src/main/target/release/main .

# Expose the port the app runs on
EXPOSE 8000

# Set the startup command to run your binary
CMD ["./main"]