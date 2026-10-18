package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsAudioError(t *testing.T) {
	cases := map[string]bool{
		"audio: audio error: device lost":       true,
		"oto: ALSA error at snd_pcm_open: busy": true,
		"graphicsdriver: failed to create GL":   false,
		"something mentioning audio: in middle": false,
	}
	for msg, want := range cases {
		if got := isAudioError(errors.New(msg)); got != want {
			t.Errorf("%q: got %v, want %v", msg, got, want)
		}
	}
	// Ebitengine wraps the backend error inside its own prefix.
	if !isAudioError(fmt.Errorf("audio: audio error: %w", errors.New("oto: x"))) {
		t.Error("wrapped audio error not detected")
	}
}
