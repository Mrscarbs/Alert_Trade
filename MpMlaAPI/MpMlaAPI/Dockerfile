# Build stage
FROM rust:latest as builder

WORKDIR /app
COPY . .

RUN apt-get update && \
    apt-get install -y pkg-config libssl-dev default-libmysqlclient-dev && \
    cargo build --release

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y libssl3 default-libmysqlclient-dev ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/MpMlaAPI /usr/local/bin/

EXPOSE 8080
CMD ["MpMlaAPI"]