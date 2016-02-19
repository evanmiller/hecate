package main

import (
//	"fmt"
	"strings"
	"github.com/nsf/termbox-go"
)

const HEAVY_VERTICAL_RIGHT = "┣"
const HEAVY_VERTICAL_LEFT = "┫"
const HEAVY_VERTICAL = "┃"
const HEAVY_HORIZONTAL = "━"
const HEAVY_TOP_LEFT = "┏"
const HEAVY_TOP_RIGHT = "┓"
const HEAVY_BOTTOM_LEFT = "┗"
const HEAVY_BOTTOM_RIGHT = "┛"

type DialogProc func(event termbox.Event, output chan<- interface{}) bool
type DialogScreen struct {
	width   int
	height  int
	top     int
	left    int
	bottom  int
	right   int
	caption string
	text    string
	callback DialogProc
}

func (screen *DialogScreen) receiveEvents (input <-chan termbox.Event, output chan<- interface{}, quit <-chan bool) {
	for {
		do_quit := false
		select {
		case event := <-input:
			do_quit = screen.callback(event, output)
		case <-quit:
			do_quit = true
		}
		if do_quit {
			break
		}
	}
}

func (screen *DialogScreen) performLayout () {
	width, height := termbox.Size()

	screen.top = Max((height - screen.height) / 2, 0)
	screen.left = Max((width - screen.width) / 2, 0)
	screen.bottom = Min(screen.top + screen.height, height)
	screen.right = Min(screen.left + screen.width, width)
}

func (screen *DialogScreen) drawScreen (style Style) {
	fg := style.default_fg
	bg := style.default_bg

	drawStringAtPoint(screen.caption, screen.left + 2, screen.top + 1, fg, bg)
	for i, s := range strings.Split(screen.text, "\n") {
		drawStringAtPoint(s, screen.left + 5, screen.top + 4 + i, fg, bg)
	}

	for row := screen.top; row < screen.bottom; row++ {
		if row == screen.bottom - 1 {
			drawStringAtPoint(HEAVY_BOTTOM_LEFT, screen.left, row, fg, bg)
			drawStringAtPoint(HEAVY_BOTTOM_RIGHT, screen.right, row, fg, bg)
		} else if row == screen.top + 2 {
			drawStringAtPoint(HEAVY_VERTICAL_RIGHT, screen.left, row, fg, bg)
			drawStringAtPoint(HEAVY_VERTICAL_LEFT, screen.right, row, fg, bg)
		} else if row == screen.top {
			drawStringAtPoint(HEAVY_TOP_LEFT, screen.left, row, fg, bg)
			drawStringAtPoint(HEAVY_TOP_RIGHT, screen.right, row, fg, bg)
		} else {
			drawStringAtPoint(HEAVY_VERTICAL, screen.left, row, fg, bg)
			drawStringAtPoint(HEAVY_VERTICAL, screen.right, row, fg, bg)
		}
		for col := screen.left + 1; col < screen.right; col++ {
			if row == screen.top || row == screen.top + 2 || row == screen.bottom - 1 {
				drawStringAtPoint(HEAVY_HORIZONTAL, col, row, fg, bg)
			}
		}
	}
}

func NewDialogScreen (caption, text string, width, height int, callback DialogProc) *DialogScreen {
	return &DialogScreen{
		text: text,
		width: width,
		height: height,
		caption: caption,
		callback: callback,
	}
}

func ShowModal (caption, text string, callback DialogProc) *DialogScreen {
	return NewDialogScreen (caption, text, 50, 8, callback)
}

func ShowMessage (caption, text string) *DialogScreen {
	return NewDialogScreen (caption, text, 70, 8, defaultDialogCloseCallback)
}

func defaultDialogCloseCallback (event termbox.Event, output chan<- interface{}) bool {
	output <- ScreenIndex(DATA_SCREEN_INDEX)
	return true
}

func Max (a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min (a, b int) int {
	if a < b {
		return a
	}
	return b
}
