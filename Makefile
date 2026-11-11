.PHONY: build run test tidy cross

build:
	go build ./...

run:
	go run ./cmd/tvshell

test:
	go test ./...

tidy:
	go mod tidy

# Cross-compile for Raspberry Pi (64-bit). Requires a C cross-toolchain
# (e.g. aarch64-linux-gnu-gcc) and CGO_ENABLED=1 for Gio's Vulkan/GL bits.
cross:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o bin/tvshell-arm64 ./cmd/tvshell
