// Retro Shooter: a top-down arcade shooter built with Ebitengine.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	err := ebiten.RunGame(game.New(scenes.NewTitle(deps)))
	if err == nil {
		return
	}
	// Ebitengine stops the game when the audio device fails, and an audio
	// context can't be created twice in one process. So keep the game
	// playable by relaunching ourselves without sound.
	if !*mute && isAudioError(err) {
		fmt.Fprintf(os.Stderr, "Sound failed (%v). Restarting without sound...\n", err)
		os.Exit(relaunchMuted())
	}
	log.Fatal(err)
}

// isAudioError reports whether err came from the sound system. Ebitengine's
// audio package prefixes its errors with "audio: ", but a device that fails
// to open returns the raw error from its backend library, oto ("oto: ").
// Neither exposes an error type, so matching prefixes is the best available.
func isAudioError(err error) bool {
	msg := err.Error()
	return strings.HasPrefix(msg, "audio: ") || strings.HasPrefix(msg, "oto: ")
}

// relaunchMuted runs this program again with -mute and returns its exit code.
func relaunchMuted() int {
	exe, err := os.Executable()
	if err != nil {
		log.Printf("can't find own executable to relaunch: %v (run with -mute)", err)
		return 1
	}
	cmd := exec.Command(exe, append([]string{"-mute"}, os.Args[1:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exit.ExitCode()
		}
		log.Printf("relaunch failed: %v (run with -mute)", err)
		return 1
	}
	return 0
}
