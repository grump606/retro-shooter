// Retro Shooter: a top-down arcade shooter built with Ebitengine.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/assets"
	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/scenes"
	"retro-shooter/internal/world"
)

func main() {
	mute := flag.Bool("mute", false, "play without sound (use if the audio device fails)")
	flag.Parse()

	ebiten.SetWindowSize(world.ScreenW*3, world.ScreenH*3)
	ebiten.SetWindowTitle("Retro Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorMode(ebiten.CursorModeHidden) // the game draws its own crosshair

	// Everything shared is built once here and handed to scenes (no globals).
	// The audio context is created exactly once: Ebitengine panics on a second.
	deps := &scenes.Deps{
		Assets: assets.Load(),
		Audio:  audio.New(*mute),
		Stars:  scenes.NewStarfield(1),
	}

	if err := ebiten.RunGame(game.New(scenes.NewTitle(deps))); err != nil {
		// Ebitengine stops the game if the audio device fails.
		if !*mute && strings.Contains(err.Error(), "audio") {
			fmt.Fprintln(os.Stderr, "Sound failed to start. Run again with -mute to play without sound.")
		}
		log.Fatal(err)
	}
}
