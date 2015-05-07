package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type PaletteScreen int

func (screen *PaletteScreen) handleKeyEvent(event termbox.Event) int {
	return DATA_SCREEN_INDEX
}

func (screen *PaletteScreen) performLayout() {
}

func (screen *PaletteScreen) drawScreen(style *Style) {
	style = style.Sub("Palette")
	width, height := termbox.Size()
	x, y := 2, 1
	for _, color := range availableColors() {
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
		x += StringOut(str, x, y, style)
		x += 2
	}
}
