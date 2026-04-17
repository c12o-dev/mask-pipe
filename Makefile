GO      ?= go
BIN     := mask-pipe
VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.PHONY: all build test clean

all: build

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/mask-pipe

test:
	$(GO) test ./...

clean:
	rm -f $(BIN)
