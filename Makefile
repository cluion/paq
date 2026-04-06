VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
BINARY  := bin/paq
LDFLAGS := -X github.com/cluion/paq/internal/cli.version=$(VERSION) \
           -X github.com/cluion/paq/internal/cli.commit=$(COMMIT) \
           -X github.com/cluion/paq/internal/cli.date=$(DATE)

.PHONY: build test lint clean install uninstall

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/paq/

test:
	go test -race -cover ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ dist/

install: build
	cp $(BINARY) /usr/local/bin/paq

uninstall:
	@rm -f /usr/local/bin/paq
	@echo "paq has been uninstalled from /usr/local/bin/paq"
