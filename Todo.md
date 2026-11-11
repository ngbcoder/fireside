# Fireside — Todo

Granular task list derived from `Design.md`. Check items off as completed; commit and push per sub-task.

Legend: `[ ]` todo · `[x]` done

---

## M0 — Skeleton
- [x] M0.1 Initialize Go module (`go.mod`) with module path (e.g. `github.com/ngbcoder/fireside`)
- [x] M0.2 Add Gio (gioui.org) dependency
- [x] M0.3 Create `cmd/tvshell/main.go` opening a Gio window
- [x] M0.4 Add `Makefile`/scripts for build and run
- [x] M0.5 Document nested `cage` dev loop (VMSVGA, Guest Additions)
- [ ] M0.6 Verify Gio window renders under nested `cage` in the VM
- [x] M0.7 Add structured logging (slog)
- [ ] M0.8 Cross-compile check: `GOOS=linux GOARCH=arm64 go build`
- [x] M0.9 Add `.gitignore` (binaries, profiles, build artifacts)

## M1 — Launcher UI
- [x] M1.1 Create `internal/ui/theme` package (colors, spacing, typography at 1080p)
- [x] M1.2 Overscan-safe margins (~5%)
- [x] M1.3 UI scale factor for 4K
- [x] M1.4 Top tab bar widget (Home, Library, App Store, Search, Settings)
- [ ] M1.5 Clock widget
- [x] M1.6 Hero shelf widget ("Continue")
- [x] M1.7 App row widget with focusable items
- [x] M1.8 "Cast to [DEVICE]" shelf widget
- [x] M1.9 Home screen composing all widgets
- [x] M1.10 Focus engine: focusable interface + focus graph
- [x] M1.11 Nearest-neighbor directional navigation
- [ ] M1.12 Select / Back / Home handling in focus engine
- [ ] M1.13 Focus visual (scale ~1.1x, shadow, label, dim unfocused)
- [x] M1.14 Keyboard input mapping to focus engine
- [ ] M1.15 Fake `InputSource` for tests
- [x] M1.16 Unit tests for focus engine navigation

## M2 — App launching
- [x] M2.1 `App` struct, `AppKind` enum (`AppWeb`, `AppFlatpak`, `AppNative`)
- [x] M2.2 `Launcher` and `Process` interfaces
- [x] M2.3 App registry loader (JSON/YAML config)
- [x] M2.4 mpv launcher implementation
- [x] M2.5 Chromium web app launcher (`--app`, Wayland backend, profile dir)
- [x] M2.6 User-agent override support
- [x] M2.7 Process tracking (`Wait`, `Stop`)
- [x] M2.8 Shell state machine (`Shell` / `App` / `Cast`)
- [ ] M2.9 Home key interception before reaching apps
- [ ] M2.10 uinput forwarder (D-pad/Select/Back as keyboard events)
- [ ] M2.11 Prototype `cage` handoff with mpv + Chromium
- [x] M2.12 Unit tests: state machine + registry parsing

## M3 — Pi bring-up
- [ ] M3.1 Flash Raspberry Pi OS Lite 64-bit; enable SSH
- [ ] M3.2 Install cage, Gio deps, mpv, Chromium on Pi
- [ ] M3.3 Cross-compile and deploy binary to Pi
- [ ] M3.4 systemd autologin + cage launch service
- [ ] M3.5 HDMI-CEC input via `/dev/cec0`
- [ ] M3.6 evdev input source (USB/BT remotes, gamepads)
- [ ] M3.7 Hardware decode test (h264_v4l2m2m etc.)
- [ ] M3.8 Boot-to-shell checklist verification
- [ ] M3.9 Thermals check under sustained playback

## M4 — AirPlay
- [x] M4.1 UxPlay supervisor (start/stop/process mgmt)
- [x] M4.2 Parse UxPlay lifecycle/session events
- [ ] M4.3 Screen handoff on `SessionStarted`, return on `SessionEnded`
- [x] M4.4 Casting overlay UI ("Casting from [DEVICE]")
- [x] M4.5 AirPlay protocol plugin implementing `Protocol` interface
- [ ] M4.6 Test cast from iPhone/iPad/Mac
- [x] M4.7 Fake `Protocol` for tests

## M5 — Web apps
- [ ] M5.1 YouTube TV web app entry (`youtube.com/tv`) + UA override
- [ ] M5.2 Disney+ web app entry
- [ ] M5.3 Apple TV web app entry (`tv.apple.com`)
- [ ] M5.4 Per-app profile directories for persistent logins
- [ ] M5.5 Verify Chromium build ships Widevine on target OS
- [ ] M5.6 Widevine quality assessment on Pi (per service)
- [ ] M5.7 Music app decision (mpv-based or web app)

## M6 — DIAL
- [x] M6.1 SSDP discovery responder (239.255.255.250:1900)
- [x] M6.2 Device description XML endpoint
- [x] M6.3 `/apps/<name>` endpoints
- [x] M6.4 YouTube launch with DIAL pairing params
- [ ] M6.5 Test with current YouTube/Netflix sender apps
- [x] M6.6 DIAL protocol plugin implementing `Protocol` interface

## M7 — App Store v1
- [x] M7.1 Curated web app registry format (JSON/YAML)
- [x] M7.2 Store backend interface
- [ ] M7.3 Featured shelf UI
- [ ] M7.4 Popular row with Get/Installed buttons
- [ ] M7.5 Category cards UI
- [ ] M7.6 Install/remove flow
- [ ] M7.7 Store states (Available/Installing/Installed/Update available)
- [ ] M7.8 Progress indicator
- [ ] M7.9 Unit tests: registry parsing

## M8 — Settings
- [x] M8.1 Settings screen scaffold
- [x] M8.2 NetworkManager D-Bus wrapper (Wi-Fi)
- [x] M8.3 BlueZ D-Bus wrapper (Bluetooth)
- [x] M8.4 Audio output selection (PipeWire)
- [ ] M8.5 Display settings
- [ ] M8.6 Device name setting
- [ ] M8.7 About screen

## M9 — Miracast (stretch)
- [ ] M9.1 Wi-Fi Direct P2P feasibility on Pi
- [ ] M9.2 Evaluate MiracleCast
- [ ] M9.3 USB Wi-Fi adapter testing (P2P-capable)
- [ ] M9.4 RTSP negotiation
- [ ] M9.5 H.264 MPEG-TS over RTP receive
- [ ] M9.6 Miracast protocol plugin implementing `Protocol` interface
- [ ] M9.7 Test with at least one sender

---

## Cross-cutting / Infrastructure
- [x] X.1 Package layout per `Design.md` §9
- [x] X.2 `internal/input` package: `Key`, `Action`, `InputEvent`, `InputSource` interface
- [x] X.3 Keyboard input source
- [x] X.4 `internal/cast/session` shared session model (`SessionKind`, `SessionEvent`, `Protocol`)
- [ ] X.5 `internal/ipc` unix socket JSON protocol
- [x] X.6 `internal/system` D-Bus wrappers (network, bluetooth, audio)
- [x] X.7 `internal/state` shell state machine
- [ ] X.8 GitHub Actions CI (build + tests)
- [ ] X.9 README with dev setup instructions

