package ui

// Layout shared by several screens, in logical pixels (320×240).
const (
	// Margin is the gap kept between screen edges and HUD/corner text.
	Margin = 4

	// Big centered banner ("WAVE 3", "PAUSED") and the line of text under it.
	BannerY     = 90
	BannerScale = 3
	BannerSubY  = BannerY + LineH*BannerScale + 2
)
