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

**Natively on Linux / WSLg** (needs a C compiler and X11/OpenGL headers):
```sh
sudo apt install gcc libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev \
  libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
go run .
```

If sound fails to start, run with `-mute`.

## Tests
```sh
go test ./internal/world ./internal/level
```
The game logic doesn't import Ebitengine, so these tests run anywhere.

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
internal/audio     audio context + synthesized sound effects
```

State flows one way: input snapshot goes down (`Game → Scene → World`), events come
back up (`World.DrainEvents` → sounds, particles, score), and scenes switch by
returning the next scene from `Update`.
