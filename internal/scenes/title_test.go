package scenes

import (
	"testing"

	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
)

func newTestTitle() *Title {
	d := &Deps{Audio: audio.New(true), Stars: NewStarfield(1)}
	s := NewTitle(d)
	// Skip the input grace period.
	for range acceptInputAfter {
		s.Update(input.State{})
	}
	return s
}

// center returns the middle of a menu item's clickable area.
func center(i menuItem) (float64, float64) {
	l, t, r, b := itemBounds(i)
	return (l + r) / 2, (t + b) / 2
}

func TestKeyboardSelectsExit(t *testing.T) {
	s := newTestTitle()
	s.Update(input.State{MenuDown: true})
	if got := s.Update(input.State{Confirm: true}); got != game.Quit {
		t.Fatalf("got %T, want game.Quit", got)
	}
}

func TestMenuWrapsAround(t *testing.T) {
	s := newTestTitle()
	s.Update(input.State{MenuUp: true}) // from START, wraps to EXIT
	if menuItem(s.selected) != itemExit {
		t.Fatalf("selected %d, want EXIT", s.selected)
	}
}

func TestClickExit(t *testing.T) {
	s := newTestTitle()
	x, y := center(itemExit)
	if got := s.Update(input.State{MouseX: x, MouseY: y, Click: true}); got != game.Quit {
		t.Fatalf("got %T, want game.Quit", got)
	}
}

func TestClickStartBeginsGame(t *testing.T) {
	s := newTestTitle()
	x, y := center(itemStart)
	if _, ok := s.Update(input.State{MouseX: x, MouseY: y, Click: true}).(*Play); !ok {
		t.Fatal("clicking START should start the game")
	}
}

func TestClickOffMenuDoesNothing(t *testing.T) {
	s := newTestTitle()
	if got := s.Update(input.State{MouseX: 5, MouseY: 5, Click: true}); got != nil {
		t.Fatalf("got %T, want nil", got)
	}
}

func TestHoverSelects(t *testing.T) {
	s := newTestTitle()
	x, y := center(itemExit)
	s.Update(input.State{MouseX: x, MouseY: y})
	if menuItem(s.selected) != itemExit {
		t.Fatal("hovering EXIT should select it")
	}
}

func TestInputIgnoredDuringGracePeriod(t *testing.T) {
	s := NewTitle(&Deps{Audio: audio.New(true), Stars: NewStarfield(1)})
	s.Update(input.State{MenuDown: true})
	if got := s.Update(input.State{Confirm: true}); got != nil {
		t.Fatalf("got %T during grace period, want nil", got)
	}
}
