package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Command struct {
	key         string
	description string
}

type AboutScreen int

func drawCommandsAtPoint(commands []Command, x int, y int, style *Style) {
	x_pos, y_pos := x, y
	for index, cmd := range commands {
		StringOut(fmt.Sprintf("%6s", cmd.key), x_pos, y_pos, style)
		StringOut(cmd.description, x_pos+8, y_pos, style)
		y_pos++
		if index > 2 && index%2 == 1 {
			y_pos++
		}
	}
}

func (screen *AboutScreen) handleKeyEvent(event termbox.Event) int {
	return DATA_SCREEN_INDEX
}

func (screen *AboutScreen) performLayout() {
}

func (screen *AboutScreen) drawScreen(style *Style) {
	style = style.Sub("About")
	logo := style.Sub("Logo")

	width, height := termbox.Size()
	template := [...]string{
		"                              ############################",
		"                              ##The#hex#editor#from#hell##",
		"                              ############################",
		"                                      ####            #   ",
		"#### #### ########  #####      ###    #### ########   #   ",
		"#### #### ####### #########  #######  #### ########  #### ",
		"#### #### ####    #### #### #### #### #### ####     ##x#x ",
		"#### #### ####    ####      #### #### #### ####       #   ",
		"######### ####### ####      ######### #### #######   ###  ",
		"######### ####### ####      ######### #### #######  # # # ",
		"#### #### ####    ####      #### #### #### ####    #  #  #",
		"#### #### ####    ####      #### #### #### ####      # #  ",
		"#### #### ####    #### #### #### #### #### ####     #   # ",
		"#### #### ####### ######### #### #### #### ####### #     #",
		"#### #### ########  #####   #### #### #### ########       ",
	}

	first_line := template[0]
	start_x := (width - len(first_line)) / 2
	start_y := (height - 2*len(template)) / 2
	x_pos := start_x
	y_pos := start_y
	for _, line := range template {
		x_pos = start_x
		for _, runeValue := range line {
			s := style

			displayRune := ' '
			if runeValue != ' ' {
				s = logo
				if runeValue != '#' {
					displayRune = runeValue
				}
				SetCell(x_pos, y_pos, displayRune, s)
			}
			x_pos++
		}
		y_pos++
	}
	commands1 := [...]Command{
		{"h", "left"},
		{"j", "down"},
		{"k", "up"},
		{"l", "right"},

		{"b", "left 4 bytes"},
		{"w", "right 4 bytes"},

		{"g", "first byte"},
		{"G", "last byte"},

		{"ctrl-e", "scroll down"},
		{"ctrl-y", "scroll up"},

		{"ctrl-f", "page down"},
		{"ctrl-b", "page up"},
	}

	commands2 := [...]Command{
		{"t", "text mode"},
		{"p", "bit pattern mode"},
		{"i", "integer mode"},
		{"f", "floating-point mode"},

		{"e", "toggle endianness"},
		{"u", "toggle signedness"},

		{"H", "shrink cursor"},
		{"L", "grow cursor"},

		{":", "jump to offset"},
		{"x", "toggle hex offset"},

		{"?", "this screen"},
		{"q", "quit program"},
	}
	x_pos = start_x + 3
	y_pos++

	drawCommandsAtPoint(commands1[:], x_pos, y_pos+1, style)
	drawCommandsAtPoint(commands2[:], x_pos+20, y_pos+1, style)
}
