VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo v1)
BINARY ?= api

build:
	go build -ldflags "-X compra-certa/api/internal/http.BuildVersion=$(VERSION)" -o $(BINARY) ./cmd/api

run:
	go run -ldflags "-X compra-certa/api/internal/http.BuildVersion=$(VERSION)" ./cmd/api
