# Fireside

A 10-foot, remote-controlled TV interface written in Go, an alternative to KDE's Plasma Bigscreen. Targets Raspberry Pi, boots into a fullscreen launcher, offers app launching, media playback, casting (AirPlay, DIAL, Miracast), and a built-in app store.

See `Design.md` for the full design doc and `Todo.md` for progress.

## Dev setup

See `docs/dev-setup.md`.

## Build

```sh
go build ./...
go test ./...
```

## Cross-compile (Pi, arm64)

Requires a C cross-toolchain (Gio uses CGO). Prefer building on the Pi or in the CachyOS VM:

```sh
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o bin/tvshell-arm64 ./cmd/tvshell
```

## Run (dev)

```sh
go run ./cmd/tvshell
# under cage:
cage -- go run ./cmd/tvshell
```
