package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Command struct {
	key         string
	description string
}

type AboutScreen struct {
	show_html bool
}

func drawCommandsAtPoint(commands []Command, x int, y int, style Style) {
	x_pos, y_pos := x, y
	longest_key_len := 1
	for _, cmd := range commands {
		if len(cmd.key) > longest_key_len {
			longest_key_len = len(cmd.key)
		}
	}
	for index, cmd := range commands {
		drawStringAtPoint(fmt.Sprintf("%[2]*[1]s", cmd.key, longest_key_len), x_pos, y_pos, style.default_fg, style.default_bg)
		drawStringAtPoint(cmd.description, x_pos+longest_key_len+2, y_pos, style.default_fg, style.default_bg)
		y_pos++
		if index > 2 && index%2 == 1 {
			y_pos++
		}
	}
}

func (screen *AboutScreen) receiveEvents(input <-chan termbox.Event, output chan<- int, quit <-chan bool) {
	for {
		do_quit := false
		select {
		case event := <-input:
			if event.Key == termbox.KeyCtrlR {
				screen.show_html = !screen.show_html
				output <- ABOUT_SCREEN_INDEX
			} else {
				screen.show_html = false
				output <- DATA_SCREEN_INDEX
			}
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

		{"^", "line start"},
		{"$", "line end"},

		{"g", "file start"},
		{"G", "file end"},

		{":", "jump to byte"},
		{"x", "toggle hex"},
	}

	commands2 := [...]Command{
		{"t", "text mode"},
		{"p", "bit pattern mode"},
		{"i", "integer mode"},
		{"f", "float mode"},

		{"H", "shrink cursor"},
		{"L", "grow cursor"},

		{"e", "toggle endianness"},
		{"u", "toggle signedness"},

		{"a", "date decoding"},
		{"@", "set date epoch"},

		{"/", "search file"},
		{"n", "next match"},
	}

	commands3 := [...]Command{
		{"S", "show tabs"},
		{"W", "hide tabs"},

		{"A", "previous tab"},
		{"D", "next tab"},

		{"ctrl-t", "new tab"},
		{"ctrl-w", "close tab"},

		{"ctrl-e", "scroll down"},
		{"ctrl-y", "scroll up"},

		{"ctrl-f", "page down"},
		{"ctrl-b", "page up"},

		{"?", "this screen"},
		{"q", "quit program"},
	}

	first_line := template[0]
	start_x := (width - len(first_line)) / 2
	start_y := (height-len(template)-2-len(commands2)/2*3-1)/2 + 1
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
	x_pos = start_x
	y_pos++

	if screen.show_html {
		drawStringAtPoint("<table>", 0, y_pos, style.default_fg, style.default_bg)
		y_pos++
		for i := 0; i < len(commands1); i++ {
			x_pos = 0
			x_pos += drawStringAtPoint("<tr>", x_pos, y_pos, style.default_fg, style.default_bg)
			for _, cmd := range [...]Command{commands1[i], commands2[i], commands3[i]} {
				x_pos += drawStringAtPoint(fmt.Sprintf("<td>%s</td>", cmd.key), x_pos, y_pos, style.default_fg, style.default_bg)
				if cmd.description == "this screen" {
					x_pos += drawStringAtPoint(fmt.Sprintf("<td>%s</td>", "help screen"), x_pos, y_pos, style.default_fg, style.default_bg)
				} else {
					x_pos += drawStringAtPoint(fmt.Sprintf("<td>%s</td>", cmd.description), x_pos, y_pos, style.default_fg, style.default_bg)
				}
			}
			x_pos += drawStringAtPoint("</tr>", x_pos, y_pos, style.default_fg, style.default_bg)
			y_pos++
		}
		drawStringAtPoint("</table>", 0, y_pos, style.default_fg, style.default_bg)
		y_pos++
	} else {
		drawCommandsAtPoint(commands1[:], x_pos, y_pos+1, style)
		drawCommandsAtPoint(commands2[:], x_pos+19, y_pos+1, style)
		drawCommandsAtPoint(commands3[:], x_pos+42, y_pos+1, style)
	}
}
