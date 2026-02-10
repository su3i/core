# Stage 1: Build
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the binary
RUN go build -o suei ./cmd/app/main.go

# Stage 2: Runtime
FROM ubuntu:22.04

# Install ffmpeg for actual transcoding
RUN apt-get update && apt-get install -y ffmpeg ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from builder stage
COPY --from=builder /app/suei /usr/local/bin/

# Set entrypoint
ENTRYPOINT ["suei"]