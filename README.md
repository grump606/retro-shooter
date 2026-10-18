# Retro Shooter

[![CI](https://github.com/grump606/retro-shooter/actions/workflows/ci.yml/badge.svg)](https://github.com/grump606/retro-shooter/actions/workflows/ci.yml)

A retro top-down arcade shooter written in Go with [Ebitengine](https://ebitengine.org).
Survive endless waves of aliens.

## Controls
| Action | Input |
|---|---|
| Move | WASD or arrow keys |
| Aim | Mouse |
| Fire | Left mouse button (hold) or Space |
| Pause | P or Esc (also pauses when the window loses focus) |
| Fullscreen | F11 or Alt+Enter |
| Debug overlay | F3 |
| Title menu | W/S or arrows to choose, Enter/Space or click to select (START / EXIT) |

## Installing on Windows
Download `RetroShooter-Setup-<version>.exe` from the latest
[CI run](https://github.com/grump606/retro-shooter/actions/workflows/ci.yml)
(artifact `retro-shooter-installer`) and run it. It installs for your user only
(no admin prompt) into `%LOCALAPPDATA%\Programs\Retro Shooter`, adds a desktop
icon (optional, on by default) and a Start menu entry, and can be removed from
Windows Settings → Apps.

Windows may warn that the installer is from an unknown publisher, because it
isn't code-signed. Choose **More info → Run anyway**.

### Building the installer
```sh
scripts/build-installer.sh   # → dist/RetroShooter-Setup-<version>.exe
```
Needs Inno Setup 6 on Windows (`winget install JRSoftware.InnoSetup --scope user`).
The version comes from the `VERSION` file. The app icon is generated from the
ship sprite (`tools/mkicon`) and embedded in the exe (`rsrc_windows_amd64.syso`);
the script regenerates both, so commit them after changing `VERSION` or the ship.

## Running

**From WSL, as a Windows program (no extra installs):**
```sh
GOOS=windows go build -o shooter.exe . && ./shooter.exe
```

**Natively on Linux / WSLg:**
```sh
go run .
```
No C compiler needed: Ebitengine loads OpenGL/X11 at runtime. Sound uses
PulseAudio (WSLg provides it) or ALSA.

If the sound device fails, the game prints the reason and restarts itself
without sound. You can also force this with `-mute`.

The repo pins Go 1.27.2 via the `toolchain` line in `go.mod`; with a newer or
older local Go, the right version is downloaded automatically.

## Adding an enemy
Every enemy's settings (stats, sprite, explosion) live in one table,
`EnemyDefs` in `internal/world/enemy.go`. Add a constant and a table entry,
then use the new kind in `internal/level/waves.go`. A test checks every entry
is complete (valid stats, equal-width sprite rows, known palette letters).

## Tests
```sh
go test ./...
```
The game logic, wave manager and screen-scaling math don't import Ebitengine,
so their tests run anywhere, with no window.

## Layout
```
main.go            window setup, builds shared deps, runs the game
VERSION            app version used by the exe and installer
installer/         Inno Setup script
scripts/           build-installer.sh
tools/mkicon       generates the app icon from the ship sprite
build/windows/     icon + exe resource config
assets/            sprites (pixel art defined in code)
internal/game      ebiten.Game: input polling, scene switching, pixel-perfect scaling
internal/scenes    title, play, game over; effects (stars, particles); rendering
internal/world     simulation: player, enemies, bullets, collisions, events
internal/level     wave definitions and spawning
internal/ui        text helpers and HUD
internal/input     per-tick input snapshot
internal/scaling   pixel-perfect scaling + mouse remapping math (tested)
internal/audio     audio context + synthesized sound effects
```

State flows one way: input snapshot goes down (`Game → Scene → World`), events come
back up (`World.DrainEvents` → sounds, particles, score), and scenes switch by
returning the next scene from `Update`.
