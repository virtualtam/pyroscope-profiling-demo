# Step 1: Build Go binaries
FROM golang:1.22-bookworm as builder

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

EXPOSE 8080

USER cook
WORKDIR /var/lib/cook

CMD ["/usr/local/bin/cook", "run"]
