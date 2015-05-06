package main

import (
	"github.com/nsf/termbox-go"
)

type Screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout()
	drawScreen(style *Style)
}

const (
	DATA_SCREEN_INDEX = iota
	ABOUT_SCREEN_INDEX
	PALETTE_SCREEN_INDEX
	EXIT_SCREEN_INDEX
)

func defaultScreensForData(bytes []byte) []Screen {
	var view_port ViewPort
	var cursor Cursor
	cursor.int_length = 4
	cursor.fp_length = 4
	cursor.mode = StringMode

	hilite := cursor.highlightRange(bytes)

	data_screen := DataScreen{bytes, cursor, hilite, view_port, cursor.mode, nil}
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

func layoutAndDrawScreen(screen Screen, style *Style) {
	screen.performLayout()
	drawBackground(style.Bg())
	screen.drawScreen(style)
	termbox.Flush()
}
