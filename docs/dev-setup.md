# Dev Setup

## Dev environment: Raspberry Pi OS Desktop in VirtualBox

The dev environment is Raspberry Pi OS (64-bit) Desktop running in VirtualBox
(previously CachyOS). This gives a Debian-based userland closer to the Pi target
while still allowing nested `cage` testing on x86_64.

- Graphics: VMSVGA, 3D acceleration on, max video memory.
- Install Guest Additions for shared clipboard/folders.
- Networking: bridged adapter, promiscuous "Allow All" for multicast discovery.

### Dev loop

```sh
# nested cage dev loop
cage -- go run ./cmd/tvshell
```

If booting cage from a TTY:
```sh
WLR_NO_HARDWARE_CURSORS=1 WLR_RENDERER_ALLOW_SOFTWARE=1 cage -- go run ./cmd/tvshell
```

## Target (Raspberry Pi)
- Raspberry Pi OS Lite or Desktop, 64-bit.
- Install: cage, mpv, chromium, libcec, build-essential.
- Build natively on the Pi or cross-compile with aarch64 toolchain.

## Cross-compile
Gio uses CGO (Vulkan/GL). Cross-compiling from a non-Pi host needs a C
cross-compiler. Prefer building on the Pi or in the RPi OS VM.
