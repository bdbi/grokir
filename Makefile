.PHONY: build clean test test-e2e lint install

BINARY=grokir
DIST=dist/$(BINARY)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.date=$(DATE)"

build:
	go build $(LDFLAGS) -o $(DIST) ./cmd/grokir

clean:
	rm -f $(DIST)

test:
	go test ./...

test-e2e:
	GROKIR_E2E=1 go test ./...

lint:
	go vet ./...

install: build
	mkdir -p ~/.local/bin
	install -m 755 $(DIST) ~/.local/bin/$(BINARY)
