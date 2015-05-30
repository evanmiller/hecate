package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout()
	drawScreen(style Style)
}

const (
	DATA_SCREEN_INDEX = iota
	ABOUT_SCREEN_INDEX
	PALETTE_SCREEN_INDEX
	EXIT_SCREEN_INDEX
)

func defaultScreensForData(bytes []byte) []Screen {
	var view_port ViewPort
	cursor := Cursor{int_length: 4, fp_length: 4, mode: StringMode,
		epoch_unit: SecondsSinceEpoch, epoch_time: time.Unix(0, 0).UTC()}

	hilite := cursor.highlightRange(bytes)

	data_screen := DataScreen{bytes: bytes, cursor: cursor, hilite: hilite, view_port: view_port, prev_mode: cursor.mode}
	about_screen := AboutScreen(0)
	palette_screen := PaletteScreen(0)
	screens := [...]Screen{
		&data_screen,
		&about_screen,
		&palette_screen,
	}

	return screens[:]
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
}

func layoutAndDrawScreen(screen Screen, style Style) {
	screen.performLayout()
	drawBackground(style.default_bg)
	screen.drawScreen(style)
	termbox.Flush()
}
