# Stage 1: Build with Go
FROM golang:1.24.5-alpine3.22 AS build

WORKDIR /build

COPY . .

RUN --mount=type=cache,id=gocache,target=/root/.cache/go-build \
    --mount=type=cache,id=gocache,target=/go/pkg/mod \
    go build -o server-executable -v main.go

# ---------------------------

# Stage 2: Minimal runtime image
FROM alpine

# Create a user and directory
RUN adduser -D -s /bin/sh -u 1000 server && \
    mkdir /server && \
    chown 1000:1000 /server

# Set working directory
WORKDIR /server

# Copy executable
COPY --chown=1000:1000 --from=build /build/server-executable .

# Copy folders (adjust as needed)
COPY --chown=1000:1000 ./config ./config
COPY --chown=1000:1000 ./maps ./maps
COPY --chown=1000:1000 ./resources ./resources
COPY --chown=1000:1000 ./translations ./translations
COPY --chown=1000:1000 ./server/worlds ./server/worlds

# Set permissions (optional if already owned)
RUN chown -R server:server /server

# Expose the port
EXPOSE 19132/udp

# Switch to the non-root user
USER server

# Start the server
CMD ["./server-executable"]
