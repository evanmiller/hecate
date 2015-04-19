package main

import (
	"github.com/nsf/termbox-go"
)

func isASCII(val byte) bool {
	return (val >= 0x20 && val < 0x7f)
}

func isCode(val byte) bool {
	return val == 0x09 || val == 0x0A || val == 0x0D
}

func isPrintable(val byte) bool {
	return isASCII(val) || isCode(val)
}

func drawStringAtPoint(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) int {
	x_pos := x
	for _, runeValue := range str {
		termbox.SetCell(x_pos, y, runeValue, fg, bg)
		x_pos++
	}
	return x_pos - x
}
