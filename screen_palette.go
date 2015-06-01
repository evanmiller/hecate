package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type PaletteScreen int

func (screen *PaletteScreen) receiveEvents(input <-chan termbox.Event, output chan<- int, quit <-chan bool) {
	for {
		do_quit := false
		select {
		case <-input:
			output <- DATA_SCREEN_INDEX
		case <-quit:
			do_quit = true
		}
		if do_quit {
			break
		}
	}
}

func (screen *PaletteScreen) performLayout() {
}

func (screen *PaletteScreen) drawScreen(style Style) {
	width, height := termbox.Size()
	fg, bg := style.default_fg, style.default_bg
	x, y := 2, 1
	for color := 1; color <= 256; color++ {
		if x+8 > width {
			x = 2
			y += 2
		}
		if y > height {
			break
		}
		termbox.SetCell(x, y, ' ', 0, termbox.Attribute(color))
		x++
		termbox.SetCell(x, y, ' ', 0, termbox.Attribute(color))
		x += 2

		str := fmt.Sprintf("%3d", color)
		x += drawStringAtPoint(str, x, y, fg, bg)
		x += 2
	}
}
