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

func drawCommandsAtPoint(commands []Command, x int, y int, style Style) {
	x_pos, y_pos := x, y
	for index, cmd := range commands {
		drawStringAtPoint(fmt.Sprintf("%6s", cmd.key), x_pos, y_pos, style.default_fg, style.default_bg)
		drawStringAtPoint(cmd.description, x_pos+8, y_pos, style.default_fg, style.default_bg)
		y_pos++
		if index > 2 && index%2 == 1 {
			y_pos++
		}
	}
}

func (screen *AboutScreen) receiveEvents(input chan termbox.Event, output chan int, quit chan bool) {
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

func (screen *AboutScreen) performLayout() {
}

func (screen *AboutScreen) drawScreen(style Style) {
	default_fg := style.default_fg
	default_bg := style.default_bg
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
	commands1 := [...]Command{
		{"h", "left"},
		{"j", "down"},
		{"k", "up"},
		{"l", "right"},

		{"b", "left 4 bytes"},
		{"w", "right 4 bytes"},

		{"g", "first byte"},
		{"G", "last byte"},

		{"/", "search file"},
		{"n", "next match"},

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

		{"D", "date decoding"},
		{"@", "set date epoch"},

		{"?", "this screen"},
		{"q", "quit program"},
	}

	first_line := template[0]
	start_x := (width - len(first_line)) / 2
	start_y := (height - len(template) - 2 - len(commands2)/2*3) / 2
	x_pos := start_x
	y_pos := start_y
	for _, line := range template {
		x_pos = start_x
		for _, runeValue := range line {
			bg := default_bg
			displayRune := ' '
			if runeValue != ' ' {
				bg = style.about_logo_bg
				if runeValue != '#' {
					displayRune = runeValue
				}
				termbox.SetCell(x_pos, y_pos, displayRune, default_fg, bg)
			}
			x_pos++
		}
		y_pos++
	}
	x_pos = start_x + 3
	y_pos++

	drawCommandsAtPoint(commands1[:], x_pos, y_pos+1, style)
	drawCommandsAtPoint(commands2[:], x_pos+20, y_pos+1, style)
}
