BUILD_DIR ?= build
SRC_FILES := $(shell find . -name "*.go")

all: lint race cover build
.PHONY: all

build: $(BUILD_DIR)/cook

$(BUILD_DIR)/%: $(SRC_FILES)
	go build -trimpath -o $@ ./cmd/$*

lint:
	golangci-lint run ./...
.PHONY: lint

cover:
	go test -coverprofile=coverage.out ./...
.PHONY: cover

coverhtml: cover
	go tool cover -html=coverage.out
.PHONY: coverhtml

race:
	go test -race ./...
.PHONY: race

test:
	go test ./...
.PHONY: test

# Live development server
live-services:
	@echo "== Starting services"
	docker compose up -d postgres pyroscope redis
.PHONY: live-services

live: live-services
	@echo "== Watching for changes... (hit Ctrl+C when done)"
	@watchexec --restart --exts go -- \
		go run ./cmd/cook/ \
			--pyroscope-addr http://localhost:4040 \
			run
.PHONY: live

live-pprof: live-services
	@echo "== Watching for changes... (hit Ctrl+C when done)"
	@watchexec --restart --exts go -- \
		go run ./cmd/cook/ \
			--pprof-addr 0.0.0.0:6060 \
			run
.PHONY: live-pprof

# pprof profiling
PPROF_URL := http://localhost:6060/debug/pprof

pprof-goroutine:
	go tool pprof $(PPROF_URL)/goroutine
.PHONY: pprof-goroutine

pprof-heap:
	go tool pprof $(PPROF_URL)/heap
.PHONY: pprof-heap

# Live development server - PostgreSQL console
psql:
	docker compose exec postgres psql -U cook restaurant
.PHONY: psql

pgdump:
	docker compose exec postgres pg_dump -U cook --format custom restaurant > restaurant.dump
.PHONY: pgdump

# Documentation
diagrams: database.png

%.png: %.d2
	d2 --sketch $< $@
