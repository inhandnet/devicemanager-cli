VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
MODULE  := github.com/inhandnet/elements-cli/internal/build

LDFLAGS := -X $(MODULE).Version=$(VERSION) \
           -X $(MODULE).Commit=$(COMMIT) \
           -X $(MODULE).Date=$(DATE)

BINARY := elements

# Detect Windows and add .exe suffix
ifeq ($(OS),Windows_NT)
  EXT := .exe
else
  EXT :=
endif

PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build build-all install clean fmt lint test docs

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)$(EXT) ./cmd/elements

build-all:
	@for pair in $(PLATFORMS); do \
		OS=$${pair%/*}; ARCH=$${pair#*/}; EXT=""; \
		[ "$$OS" = "windows" ] && EXT=".exe"; \
		echo "Building $${OS}/$${ARCH}..."; \
		CGO_ENABLED=0 GOOS=$$OS GOARCH=$$ARCH go build \
			-ldflags "$(LDFLAGS)" \
			-o bin/$(BINARY)-$${OS}-$${ARCH}$${EXT} ./cmd/elements; \
	done

install:
	CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" ./cmd/elements

test:
	CGO_ENABLED=0 go test ./... -v

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

docs:
	go run ./cmd/docgen
