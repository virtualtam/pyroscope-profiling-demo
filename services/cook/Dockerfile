# Step 1: Build Go binaries
FROM golang:1.24-bookworm AS builder

ARG CGO_ENABLED=1

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

ADD . .
RUN --mount=type=cache,target=/root/.cache/go-build make build

# Step 2: Build the actual image
FROM debian:bookworm-slim

RUN groupadd \
    --gid 1000 \
    cook \
    && useradd \
    --create-home \
    --home-dir /var/lib/cook \
    --shell /bin/bash \
    --uid 1000 \
    --gid cook \
    cook

COPY --from=builder /app/build/cook /usr/local/bin/cook

ENV COOK_LISTEN_ADDR="0.0.0.0:8080" \
    COOK_LOG_FORMAT="json" \
    COOK_LOG_LEVEL="info" \
    COOK_DB_HOST="postgres" \
    COOK_DB_PORT=5432 \
    COOK_REDIS_HOST="redis" \
    COOK_REDIS_PORT=6379 \
    COOK_REDIS_DB=0 \
    COOK_PPROF_ADDR="" \
    COOK_PYROSCOPE_ADDR="http://pyroscope:4040"

# pprof
EXPOSE 6060

# cook API
EXPOSE 8080

USER cook
WORKDIR /var/lib/cook

CMD ["/usr/local/bin/cook", "run"]
