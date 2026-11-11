# Fireside: Design Doc

**Status:** Draft v0.1 (idea / design stage)
**Name:** Fireside
**Mockup:** https://claude.ai/artifact/S6gtyy7u3mnzshyZ5j5ZiK (Home and App Store screens)

---

## 1. Overview

A 10-foot, remote-controlled TV interface written in Go, intended as an alternative to KDE's Plasma Bigscreen. It runs on a Raspberry Pi, boots straight into a fullscreen launcher, and offers app launching, media playback, casting (AirPlay, DIAL, Miracast), and a built-in app store. The look is inspired by Apple TV: a top tab bar, a hero shelf, and focus-driven icon rows over a soft background.

### Goals

- Boot to a polished, remote-friendly launcher on a Raspberry Pi.
- Navigate entirely with a TV remote (HDMI-CEC), USB remote, gamepad, or keyboard.
- Receive casts from phones and computers via AirPlay, DIAL, and Miracast.
- Run streaming services (YouTube, Disney+, Apple TV) as first-class "apps".
- Include an app store for installing more apps.
- Keep the shell small: mostly Go, few moving parts, easy to cross-compile.

### Non-goals (for now)

- Writing a custom Wayland compositor.
- Official Chromecast ("Google Cast") receiver support (requires Google licensing; see section 6).
- Building native clients for DRM-protected services.
- Replacing a full desktop environment.

---

## 2. Target Platforms and Dev Environment

| Item | Choice |
|---|---|
| Primary hardware | Raspberry Pi (Pi 4 or Pi 5 preferred; Pi 3 is possible but tight) |
| OS base | Raspberry Pi OS Lite, 64-bit (no desktop), or similar minimal Linux |
| Dev environment | CachyOS (Arch-based) VM in VirtualBox |
| Build | Go, native on x86_64 for dev; `GOOS=linux GOARCH=arm64` or native build on the Pi for release |

### What the VM can and cannot test

| Feature | In the VM? |
|---|---|
| Gio UI, focus navigation, animations | Yes (enable 3D acceleration, VMSVGA) |
| Keyboard as a stand-in remote | Yes |
| mpv playback | Yes (software decode only) |
| App launching and process tracking | Yes |
| DIAL and AirPlay discovery | Yes, with a **bridged** network adapter |
| HDMI-CEC | No |
| Hardware video decode | No (Pi-specific) |
| Miracast / Wi-Fi Direct | Essentially no |
| Widevine playback quality | Verify on the Pi |

### VirtualBox notes

- Graphics controller: VMSVGA, 3D acceleration on, video memory at maximum.
- Install Guest Additions (`virtualbox-guest-utils` on Arch) and enable its service.
- Prefer running `cage` nested inside the VM's desktop session as the dev loop. Booting `cage` directly from a TTY may need `WLR_NO_HARDWARE_CURSORS=1` and `WLR_RENDERER_ALLOW_SOFTWARE=1`.
- Networking: bridged adapter, promiscuous mode "Allow All" if discovery is flaky. Multicast over a Wi-Fi host bridge can be unreliable; a wired host is more dependable.

---

## 3. Architecture

```
+----------------------------------------------------------+
|                    cage (Wayland kiosk)                  |
|  +----------------------------------------------------+  |
|  |  tvshell (Go, Gio)                                 |  |
|  |   UI  |  Focus engine  |  App manager  |  Store    |  |
|  +-------------------+--------------------------------+  |
|                      | local IPC (unix socket, JSON)     |
|  +-------------------+--------------------------------+  |
|  |  castd (Go)       |   inputd (Go)                  |  |
|  |  AirPlay / DIAL / |   CEC / evdev / keyboard ->    |  |
|  |  Miracast plugins |   normalized key events        |  |
|  +-------------------+--------------------------------+  |
|                                                          |
|  External processes: mpv, UxPlay, Chromium, flatpak apps |
+----------------------------------------------------------+
   System services: NetworkManager, BlueZ, PipeWire, D-Bus
```

### Components

- **tvshell:** the launcher UI. Renders with Gio, owns the focus model and navigation, shows the home screen, store, settings, and casting overlays.
- **inputd:** reads all input devices and emits normalized key events. Also injects keyboard events into other apps through `/dev/uinput`.
- **castd:** advertises the TV on the network and manages incoming cast sessions through protocol plugins.
- **App manager:** launches, tracks, and stops external apps (web apps, Flatpaks, native binaries).
- **Store:** browses and installs apps from a backend (see section 8).
- **System layer:** wraps NetworkManager, BlueZ, and PipeWire over D-Bus (`godbus`).

Each daemon can start as a goroutine inside one binary, then be split into separate processes only if isolation or restart behavior needs it.

### Screen ownership

Only one thing owns the screen at a time: the shell, a launched app, or a cast session. The shell keeps a small state machine (`Shell`, `App`, `Cast`) and hands focus over when a session starts and takes it back when it ends. A dedicated Home key is handled by `inputd` before it reaches any app, so a stuck app can always be dismissed.

> **Open question:** `cage` is designed to show a single application fullscreen. Verify how it behaves when the shell launches additional clients (Chromium, UxPlay, mpv). If it doesn't layer or switch the way we need, evaluate `labwc` or a kiosk-configured `sway` instead.

---

## 4. UI Design

Inspired by Apple TV. It does not copy Apple's artwork, icons, or branding.

### Visual language

- Dark base with a soft blurred wallpaper (generated; user-replaceable).
- Frosted-glass panels, large rounded corners, generous spacing.
- Large type: 72 px hero titles, 30 px section headings, 24-28 px body, all at 1920x1080.
- Focus: the focused item scales up (about 1.1x), lifts with a shadow, and shows its label; unfocused items dim slightly.
- Line icons throughout (generic symbols, replaced per app with real app icons later).

### Screens

1. **Home:** top tab bar (Home, Library, App Store, Search, Settings), clock, hero shelf ("Continue"), app row, and a "Cast to [DEVICE NAME]" shelf.
2. **App Store:** featured shelf, Popular row with Get/Installed buttons, category cards.
3. **Settings** (not yet mocked): network, Bluetooth, audio output, display, device name, about.
4. **Casting overlay** (not yet mocked): "Casting from [DEVICE]" with pause/stop, and a full takeover while video plays.
5. **Library** and **Search** (not yet mocked).

### Navigation model

A 2D focus graph. Directions move focus to the nearest focusable neighbor; Select activates; Back returns up one level; Home returns to the launcher from anywhere. Every screen must be fully usable with only these keys.

### Accessibility and TV constraints

- Text no smaller than 22 px at 1080p; maintain strong contrast on glass panels.
- Overscan-safe margins (keep important content at least about 5% from edges).
- Focus is always visible; no hover-only affordances.
- Support 1080p first; scale for 4K through a UI scale factor.

---

## 5. Input

All devices are normalized behind one interface.

```go
type Key int

const (
	KeyUp Key = iota
	KeyDown
	KeyLeft
	KeyRight
	KeySelect
	KeyBack
	KeyHome
	KeyPlayPause
	KeyMenu
	KeyVolumeUp
	KeyVolumeDown
)

type Action int

const (
	Press Action = iota
	Release
	Repeat
)

type InputEvent struct {
	Key    Key
	Action Action
	Source string // "cec", "evdev:/dev/input/event3", "keyboard"
}

type InputSource interface {
	Name() string
	Events() <-chan InputEvent
	Close() error
}
```

### Sources

| Source | Notes |
|---|---|
| Keyboard | First implementation; works in the VM |
| HDMI-CEC | Kernel CEC API (`/dev/cec0`) on the Pi, or `libcec` |
| evdev | USB/Bluetooth remotes, gamepads, air mice |
| Virtual (uinput) | Output side: translates shell-level keys into arrow/Enter/Escape for web apps |

### Key translation for web apps

When a web app is in the foreground, `inputd` forwards D-Pad, Select, and Back as standard keyboard events through uinput, while still intercepting Home (and optionally a long-press Back) for the shell.

---

## 6. Casting

Casting runs in `castd`. Each protocol is a plugin feeding a common session model, so the shell does not care which protocol started a session.

### Why not "Chromecast"

A Google Cast receiver is not an open standard. Senders discover receivers over mDNS and talk CASTV2 on port 8009, and licensed receivers pass a device-authentication step using certificates issued by Google. Unlicensed receivers may be hidden or refused by many sender apps, and official "Chromecast built-in" requires Google's licensing program. This project therefore targets the three open or reverse-engineered protocols below. An experimental CASTV2 receiver can be revisited later; its device-auth behavior should be checked against current sender apps before investing.

### Protocol plan (build in this order)

**1. AirPlay (first)**

- Wrap **UxPlay** (open source, GStreamer-based) as a managed subprocess.
- Covers iPhone, iPad, and Mac: mirroring, video, and audio.
- UxPlay opens its own video window, so the shell hands over the screen and takes it back when the process ends.
- Do not reimplement AirPlay in Go.

**2. DIAL**

- SSDP discovery on 239.255.255.250:1900, plus a small HTTP service (device description XML and `/apps/<name>` endpoints).
- For YouTube, the usual approach is to launch the YouTube TV web app in the kiosk browser with the pairing parameters DIAL provides.
- Sender-app support for DIAL has been shrinking and some services have cut back casting. Test the current YouTube and Netflix apps early before depending on it. Expect YouTube to be the main win.
- Go is a good fit: it is just UDP multicast and HTTP.

**3. Miracast (last, stretch goal)**

- Wi-Fi Direct (P2P) link, then RTSP negotiation and H.264 MPEG-TS over RTP.
- Linux receiver support is weak (MiracleCast is largely stagnant).
- The Pi's onboard Wi-Fi may not support P2P reliably, and one radio struggles to hold home Wi-Fi and Wi-Fi Direct simultaneously. A second USB Wi-Fi adapter with solid P2P support is the realistic fix.
- Not all senders support it (many Pixel phones do not).
- Do not block any other milestone on Miracast.

### castd interface

```go
type SessionKind int

const (
	SessionStarted SessionKind = iota
	SessionMedia      // media loaded or changed
	SessionPaused
	SessionResumed
	SessionEnded
)

type SessionEvent struct {
	Kind     SessionKind
	Protocol string // "airplay", "dial", "miracast"
	Peer     string // sender device name, if known
	Title    string // optional, if the protocol provides it
}

type Protocol interface {
	Name() string
	Start(ctx context.Context, deviceName string) error // begin advertising
	Stop() error
	Events() <-chan SessionEvent
}
```

The shell subscribes to events from all protocols, switches its state machine to `Cast` on `SessionStarted`, and returns to `Shell` on `SessionEnded`.

---

## 7. Apps

An "app" is anything the shell can launch.

```go
type AppKind int

const (
	AppWeb AppKind = iota
	AppFlatpak
	AppNative
)

type App struct {
	ID         string
	Name       string
	Kind       AppKind
	URL        string   // web apps
	Exec       string   // native apps / flatpak ref
	Icon       string
	UserAgent  string   // optional override (e.g. a TV user agent)
	ProfileDir string   // per-app browser profile for persistent logins
}

type Launcher interface {
	Launch(ctx context.Context, app App) (Process, error)
}

type Process interface {
	Wait() error
	Stop() error
}
```

### Built-in apps

| App | Approach | Notes |
|---|---|---|
| YouTube | Web app: YouTube TV (`youtube.com/tv`) in kiosk Chromium; optional native client later | The TV site has its own remote-friendly UI and works with DIAL |
| Disney+ | Web app wrapper only | DRM-protected, no public API |
| Apple TV | Web app wrapper only (`tv.apple.com`) | DRM-protected, no public API |
| Media Player | mpv | Local files and network streams |
| Music | TBD (mpv-based or web app) | |
| Settings | Native (shell screen) | |

### Web app launching

- Start Chromium in app/kiosk mode (`--app=<url>`) on the Wayland backend under `cage`.
- Use a **separate profile directory per app** so logins persist and don't mix.
- Apply an optional user-agent override so services serve their TV layouts.
- Forward remote input as keyboard events (section 5); keep Home reserved for the shell.

### DRM and playback quality

- Disney+ and Apple TV protect video with DRM (Widevine in Chromium). On Linux, browsers generally get only the lowest Widevine security level, which often limits video to standard definition or around 720p rather than 1080p or 4K. Limits differ per service and change over time.
- Confirm that the Chromium build on the target OS ships with Widevine; some distro builds do not and it may need to be added.
- **Test on the Pi early.** The VM cannot answer this question.

### Optional: native YouTube client

A custom YouTube UI could sit on top of an open front end such as Invidious or Piped, which avoids DRM. This sits in a gray area of YouTube's terms and can break when YouTube changes. Treat it as an experiment, separate from the main YouTube web app.

---

## 8. App Store

The store UI is mocked up; the backend is **undecided**.

| Option | Pros | Cons |
|---|---|---|
| **Flatpak / Flathub** | Large catalog, handles install/update/sandboxing, minimal custom work | Most apps are desktop-oriented, not TV-friendly; needs curation and a TV-usable filter |
| **Curated web app registry** | Simple, TV-friendly, easy to add services as a JSON entry | Small catalog; relies on each site working well with a remote |
| **Own repository** | Full control, TV-focused apps only | Most work; you host and maintain everything |

A reasonable path: start with a **curated web app registry** (JSON/YAML entries that become `App` records), add **Flatpak** for a vetted subset, and revisit a custom repository only if needed.

### Store states

Available (Get), Installing (progress), Installed (Open/Remove), Update available.

---

## 9. Suggested Package Layout

```
tvshell/
  cmd/
    tvshell/            # main launcher binary
    castd/              # optional standalone cast daemon
  internal/
    ui/                 # Gio screens, widgets, theme, focus engine
      home/
      store/
      settings/
      cast/
    input/              # InputSource implementations + uinput forwarder
      keyboard/
      cec/
      evdev/
    cast/
      airplay/          # UxPlay supervisor
      dial/             # SSDP + HTTP
      miracast/         # experimental
      session/          # shared session model
    apps/               # App registry, launchers (web, flatpak, native)
    store/              # backends (registry, flatpak)
    system/             # D-Bus wrappers: network, bluetooth, audio
    ipc/                # unix socket protocol between components
    state/              # shell state machine (Shell / App / Cast)
  assets/               # wallpaper, icons, fonts
  configs/              # default apps, device name, theme
  docs/
```

---

## 10. Milestones

| # | Milestone | Done when |
|---|---|---|
| M0 | Skeleton | Go module, Gio window under nested `cage` in the VM |
| M1 | Launcher UI | Home screen with tab bar, hero, app row, focus navigation with keyboard |
| M2 | App launching | Launch and track mpv and a web app; Home key returns to shell |
| M3 | Pi bring-up | Boots to the shell on a Pi; HDMI-CEC input working; mpv playback tested |
| M4 | AirPlay | UxPlay supervised; cast from an iPhone/Mac takes over the screen and returns |
| M5 | Web apps | YouTube, Disney+, Apple TV launch with persistent logins; Widevine quality assessed on the Pi |
| M6 | DIAL | Discoverable; YouTube launch from a phone works |
| M7 | App Store v1 | Curated web app registry, install/remove from the UI |
| M8 | Settings | Wi-Fi, Bluetooth, audio output, device name |
| M9 | Miracast (stretch) | Works with at least one sender on supported hardware |

---

## 11. Risks and Open Questions

| Risk / question | Impact | Mitigation |
|---|---|---|
| `cage` behavior with multiple clients | Core launch/handoff flow | Prototype early in M2; fall back to `labwc` or kiosk `sway` |
| Widevine availability and quality on the Pi | Disney+/Apple TV usefulness | Test in M5; set expectations; document limits |
| DIAL support shrinking in sender apps | Casting reach | Test current apps early; treat as best-effort |
| Miracast on Pi Wi-Fi hardware | Feature may not be viable | Stretch goal; optional USB adapter |
| Gio C dependencies and cross-compiling | Build friction | Build on the Pi or use a cross toolchain/Docker |
| Hardware video decode differences (Pi 4 vs Pi 5) | Playback performance | Test both; document supported models |
| Web apps and remote navigation | Some sites are hard to drive with a D-Pad | Prefer TV-oriented sites; consider a cursor-mode fallback |
| Terms of service for wrapped services and unofficial clients | Legal/ToS risk | Prefer official web apps; keep unofficial clients experimental |
| Store backend choice | Scope | Start with curated registry; revisit later |
| Project name and branding | Cosmetic | Decide before public release |

---

## 12. Testing Strategy

- **Unit tests:** focus engine (graph navigation), input normalization, app registry parsing, session state machine.
- **Interface fakes:** a fake `InputSource` and a fake `Protocol` to drive the UI without hardware.
- **VM integration:** nested `cage` run, mpv playback, app launch/stop, DIAL and AirPlay discovery (bridged network).
- **Pi hardware checklist:** boot-to-shell, CEC remote, hardware decode, Widevine playback per service, audio output, cast sessions with real phones, thermals under sustained playback.

---

## 13. Next Steps

1. Create the Go module and a Gio window running under nested `cage` (M0).
2. Build the Home screen from the mockup with keyboard navigation (M1).
3. Prototype launching mpv and a Chromium web app and verify the `cage` handoff (M2).
4. Order or set up a Pi and verify CEC and Widevine early, since they carry the biggest unknowns.
