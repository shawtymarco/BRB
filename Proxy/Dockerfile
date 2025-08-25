# Stage 1: Build with Go
FROM golang:1.25.0-alpine3.22 AS build

WORKDIR /app

# Copy module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code maintaining directory structure
COPY . .

# Build from the root of the module
RUN --mount=type=cache,id=gocache,target=/root/.cache/go-build \
    --mount=type=cache,id=gocache,target=/go/pkg/mod \
    go build -o proxy-executable -v ./main.go

# ---------------------------
# Stage 2 remains the same
FROM alpine

RUN adduser -D -s /bin/sh -u 1000 proxy && \
    mkdir /proxy && \
    chown 1000:1000 /proxy

WORKDIR /proxy

COPY --chown=1000:1000 --from=build /app/proxy-executable .
COPY --chown=1000:1000 config.toml .

RUN chown -R proxy:proxy /proxy

EXPOSE 19132/udp

USER proxy

CMD ["./proxy-executable"]