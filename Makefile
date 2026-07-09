GO_IMAGE ?= golang:1.26.5-alpine
GO := go

ifneq ($(shell uname -s),Linux)
GO := docker run --rm \
	--user "$$(id -u):$$(id -g)" \
	-e HOME=/tmp \
	-e GOCACHE=/tmp/go-build \
	-e GOMODCACHE=/tmp/go-mod \
	-v "$(CURDIR):/workspace" \
	-w /workspace \
	$(GO_IMAGE) go
endif

.PHONY: build test lint fmt clean

build:
	mkdir -p build
	$(GO) build -trimpath -o ./build/home-backup ./cmd/home-backup

test:
	$(GO) test ./...

lint:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf build
