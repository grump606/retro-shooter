# Retro Shooter

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

## Tests
```sh
go test ./...
```
The game logic, wave manager and screen-scaling math don't import Ebitengine,
so their tests run anywhere, with no window.

## Layout
```
main.go            window setup, builds shared deps, runs the game
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
