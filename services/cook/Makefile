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
