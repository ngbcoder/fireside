# Dev Setup

## VM (CachyOS in VirtualBox)
- Graphics: VMSVGA, 3D acceleration on, max video memory.
- Install Guest Additions: `virtualbox-guest-utils`, enable service.
- Networking: bridged adapter, promiscuous "Allow All" for multicast discovery.
- Dev loop: run `cage` nested inside the desktop session.

```sh
# nested cage dev loop
cage -- go run ./cmd/tvshell
```

If booting cage from a TTY:
```sh
WLR_NO_HARDWARE_CURSORS=1 WLR_RENDERER_ALLOW_SOFTWARE=1 cage -- go run ./cmd/tvshell
```

## Pi (target)
- Raspberry Pi OS Lite, 64-bit.
- Install: cage, mpv, chromium, libcec, build-essential.
- Build natively on the Pi or cross-compile with aarch64 toolchain.

## Cross-compile
Gio uses CGO (Vulkan/GL). Cross-compiling from a non-Linux host needs a C
cross-compiler. Prefer building on the Pi or in the CachyOS VM.
