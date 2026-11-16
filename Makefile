GO      ?= go
BIN_DIR ?= bin
ADDR    ?= 127.0.0.1:4242

SERVER_BIN := $(BIN_DIR)/tap-server
CLI_BIN    := $(BIN_DIR)/tap-cli
GUI_BIN    := $(BIN_DIR)/tap-gui

.DEFAULT_GOAL := all

all: build

help:
	@printf '%s\n' 'make: build the server and CLI' 'make build-gui: build the Fyne GUI' 'make run-server / run-client / run-client-gui' 'make test: basic server and CLI tests'

install:
	$(GO) mod download

build: build-server build-client

build-server:
	$(GO) build -buildvcs=false -o $(SERVER_BIN) ./cmd/server

build-client:
	$(GO) build -buildvcs=false -o $(CLI_BIN) ./cmd/cli

build-gui:
	$(GO) build -buildvcs=false -o $(GUI_BIN) ./cmd/gui

run-server:
	$(GO) run -buildvcs=false ./cmd/server

run-client:
	$(GO) run -buildvcs=false ./cmd/cli "$(ADDR)"

run-client-gui:
	LANGUAGE=$$(printf '%s' "$$LANGUAGE" | sed 's/^:*//;s/:*$$//;s/::*/:/g') $(GO) run -buildvcs=false ./cmd/gui

test:
	$(GO) test -timeout 30s ./cmd/server ./cmd/cli
	$(GO) test -tags ci -timeout 30s ./cmd/gui

lint:
	@unformatted="$$(gofmt -l cmd)"; \
	if [ -n "$$unformatted" ]; then printf 'gofmt needed for:\n%s\n' "$$unformatted"; exit 1; fi
	$(GO) vet ./cmd/server ./cmd/cli

lint-gui:
	$(GO) vet ./cmd/gui

clean:
	rm -rf $(BIN_DIR)

.PHONY: all help install build build-server build-client build-gui \
        run-server run-client run-client-gui test lint lint-gui clean
