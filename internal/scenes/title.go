package scenes

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
	"retro-shooter/internal/ui"
	"retro-shooter/internal/world"
)

// Title is the start screen with a START / EXIT menu.
//
// Requiring a click or key press here also unlocks audio: browsers (a future
// web build) refuse to play sound until the user interacts with the page, and
// Ebitengine resumes its audio context on that first input. By the time
// gameplay starts, sound is guaranteed to work.
type Title struct {
	d        *Deps
	t        int
	selected int
	lastMX   float64
	lastMY   float64
}

type menuItem int

const (
	itemStart menuItem = iota
	itemExit
	numItems
)

var menuLabels = [numItems]string{itemStart: "START", itemExit: "EXIT"}

// Title screen layout in logical pixels.
const (
	titleY        = 30 // "RETRO"
	titleScale    = 4
	subtitleY     = 90 // "SHOOTER"
	subtitleScale = 3
	controlsY     = 196 // two lines of control hints near the bottom

	menuTop     = 140 // y of the first item
	menuSpacing = 18  // distance between items
	menuPadX    = 12  // clickable area extends this far beyond the label
	menuBoxGap  = 2   // the highlight box stops this short of the next item
)

// acceptInputAfter ignores input for a few ticks so a click or key press from
// the previous screen doesn't immediately pick a menu item.
const acceptInputAfter = 15

func NewTitle(d *Deps) *Title { return &Title{d: d, lastMX: -1, lastMY: -1} }

func (s *Title) Update(in input.State) game.Scene {
	s.d.Audio.NextTick()
	s.d.Stars.Update()
	s.t++

	// Hovering selects an item, but only when the mouse actually moves, so a
	// resting cursor doesn't fight the keyboard.
	mouseMoved := in.MouseX != s.lastMX || in.MouseY != s.lastMY
	s.lastMX, s.lastMY = in.MouseX, in.MouseY
	hovered, onItem := itemAt(in.MouseX, in.MouseY)
	if mouseMoved && onItem && hovered != menuItem(s.selected) {
		s.selected = int(hovered)
		s.d.Audio.Play(audio.Blip)
	}

	if s.t <= acceptInputAfter {
		return nil
	}

	if in.MenuUp {
		s.selected = (s.selected + int(numItems) - 1) % int(numItems)
		s.d.Audio.Play(audio.Blip)
	}
	if in.MenuDown {
		s.selected = (s.selected + 1) % int(numItems)
		s.d.Audio.Play(audio.Blip)
	}

	switch {
	case in.Click && onItem:
		return s.activate(hovered)
	case in.Confirm:
		return s.activate(menuItem(s.selected))
	}
	return nil
}

func (s *Title) activate(item menuItem) game.Scene {
	switch item {
	case itemExit:
		return game.Quit
	default:
		s.d.Audio.Play(audio.Blip)
		return NewPlay(s.d)
	}
}

// itemAt returns the menu item under the given logical point, if any.
func itemAt(x, y float64) (menuItem, bool) {
	for i := range numItems {
		left, top, right, bottom := itemBounds(i)
		if x >= left && x < right && y >= top && y < bottom {
			return i, true
		}
	}
	return 0, false
}

func itemBounds(i menuItem) (left, top, right, bottom float64) {
	w := float64(len(menuLabels[i]) * ui.CharW)
	left = (world.ScreenW-w)/2 - menuPadX
	top = float64(menuTop + int(i)*menuSpacing)
	return left, top, left + w + 2*menuPadX, top + menuSpacing
}

func (s *Title) Draw(screen *ebiten.Image) {
	screen.Fill(ui.Night)
	s.d.Stars.Draw(screen, s.d.Assets.Pixel)

	ui.BigText(screen, "RETRO", titleY, titleScale, ui.Cyan)
	ui.BigText(screen, "SHOOTER", subtitleY, subtitleScale, ui.Yellow)

	for i := range numItems {
		y := menuTop + int(i)*menuSpacing
		if int(i) == s.selected {
			left, top, right, bottom := itemBounds(i)
			ui.Rect(screen, left, top, right-left, bottom-top-menuBoxGap, ui.Dark)
			ui.BigText(screen, "> "+menuLabels[i]+" <", y, 1, ui.Yellow)
		} else {
			ui.BigText(screen, menuLabels[i], y, 1, ui.Grey)
		}
	}

	ui.TextCentered(screen, "WASD MOVE  MOUSE AIM  CLICK FIRE", controlsY)
	ui.TextCentered(screen, "P PAUSE   F11 FULLSCREEN", controlsY+ui.LineH)
	if s.d.HighScore > 0 {
		ui.TextCentered(screen, fmt.Sprintf("HI %06d", s.d.HighScore), ui.Margin)
	}
}
