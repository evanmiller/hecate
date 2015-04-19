package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DisplayScreen int

const (
	DataScreen = iota
	ColorScreen
	AboutScreen
)

type ViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type Style struct {
	default_bg termbox.Attribute
	default_fg termbox.Attribute

	rune_fg       termbox.Attribute
	space_rune_fg termbox.Attribute
	int_fg        termbox.Attribute
	bit_fg        termbox.Attribute

	text_cursor_hex_bg termbox.Attribute
	bit_cursor_hex_bg  termbox.Attribute
	int_cursor_hex_bg  termbox.Attribute
	fp_cursor_hex_bg   termbox.Attribute

	hilite_hex_fg  termbox.Attribute
	hilite_rune_fg termbox.Attribute
}

func drawAboutScreen(default_fg termbox.Attribute, default_bg termbox.Attribute) {
	drawBackground(default_bg)
	width, height := termbox.Size()
	/* Well, this is awfully dark! */
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
	start_y := (height - len(template)) / 2
	x_pos := start_x
	y_pos := start_y
	for _, line := range template {
		x_pos = start_x
		for _, runeValue := range line {
			bg := default_bg
			displayRune := ' '
			if runeValue != ' ' {
				bg = termbox.Attribute(125)
				if runeValue != '#' {
					displayRune = runeValue
				}
			}
			termbox.SetCell(x_pos, y_pos, displayRune, default_fg, bg)
			x_pos++
		}
		y_pos++
	}
	termbox.Flush()
}

func isASCII(val byte) bool {
	return (val >= 0x20 && val < 0x80)
}

func isCode(val byte) bool {
	return val == 0x09 || val == 0x0A || val == 0x0D
}

func isPrintable(val byte) bool {
	return isASCII(val) || isCode(val)
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
	/*
		width, height := termbox.Size()
		x, y := 0, 0
		for x = 0; x < width; x++ {
			for y = 0; y < height; y++ {
				termbox.SetCell(x, y, ' ', 0, bg)
			}
		} */
}

func drawStringAtPoint(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) int {
	x_pos := x
	for _, runeValue := range str {
		termbox.SetCell(x_pos, y, runeValue, fg, bg)
		x_pos++
	}
	return x_pos - x
}

func drawDataScreen(data []byte, old_view_port ViewPort, style Style, cursor Cursor, hilite ByteRange) ViewPort {
	x, y := 2, 1
	width, height := termbox.Size()
	rows := 1
	drawBackground(style.default_bg)
	legend_height := drawWidgets(cursor, style)

	var new_view_port ViewPort
	new_view_port.bytes_per_row = (width - 3) / 3
	new_view_port.number_of_rows = (height - 1 - legend_height) / 3

	cursor_row_within_view_port := 0
	if old_view_port.bytes_per_row > 0 {
		cursor_row_within_view_port = cursor.pos/old_view_port.bytes_per_row - old_view_port.first_row
		if cursor.pos/new_view_port.bytes_per_row > cursor_row_within_view_port {
			new_view_port.first_row = cursor.pos/new_view_port.bytes_per_row - cursor_row_within_view_port
		}
		if cursor.pos/new_view_port.bytes_per_row >= new_view_port.first_row+new_view_port.number_of_rows {
			new_view_port.first_row = cursor.pos/new_view_port.bytes_per_row - new_view_port.number_of_rows + 1
		}
	}

	start := new_view_port.first_row * new_view_port.bytes_per_row
	end := start + new_view_port.number_of_rows*new_view_port.bytes_per_row
	for index := start; index < end && index < len(data); index++ {
		b := data[index]
		hex_fg := style.default_fg
		hex_bg := style.default_bg
		code_fg := style.space_rune_fg
		rune_fg := style.rune_fg
		rune_bg := style.default_bg
		cursor_length := cursor.length()
		if x+3 > width-1 {
			x = 2
			y += 3
			rows++
		}
		if y > height-2 {
			break
		}
		if index >= cursor.pos && index < cursor.pos+cursor_length {
			hex_bg = cursor.color(style)
			termbox.SetCell(x-1, y, ' ', hex_fg, hex_bg)
			termbox.SetCell(x+2, y, ' ', hex_fg, hex_bg)
		} else if index >= hilite.pos && index < hilite.pos+hilite.length {
			hex_fg = style.hilite_hex_fg
		}
		if index >= hilite.pos && index < hilite.pos+hilite.length {
			rune_fg = style.hilite_rune_fg
			code_fg = style.rune_fg
		}
		if cursor.mode == StringMode || index < cursor.pos || index >= cursor.pos+cursor_length {
			if b == 0x20 {
				termbox.SetCell(x, y+1, '•', style.space_rune_fg, rune_bg)
			} else if isASCII(b) {
				termbox.SetCell(x, y+1, rune(b), rune_fg, rune_bg)
			} else if isCode(b) {
				codes := map[byte]rune{
					0x0A: 'n',
					0x0D: 'r',
					0x09: 't',
				}
				termbox.SetCell(x, y+1, '\\', code_fg, rune_bg)
				termbox.SetCell(x+1, y+1, codes[b], code_fg, rune_bg)
			} else {
				termbox.SetCell(x, y+1, ' ', 0, rune_bg)
			}
		} else if cursor.mode == BitPatternMode {
			for i := 0; i < 8; i++ {
				if b&(1<<uint8(7-i)) > 0 {
					termbox.SetCell(x-1+(i%4), y+1+i/4, '●', style.bit_fg, rune_bg)
				} else {
					termbox.SetCell(x-1+(i%4), y+1+i/4, '○', style.bit_fg, rune_bg)
				}
			}
		} else if index == cursor.pos {
			cursor_length := cursor.length()
			total_length := cursor_length*3 + 1
			str := cursor.formatBytesAsNumber(data[cursor.pos : cursor.pos+cursor_length])
			x_copy := x - 1
			y_copy := y + 1
			x_copy = x_copy + (total_length-len(str))/2
			for _, runeValue := range str {
				termbox.SetCell(x_copy, y_copy, runeValue, style.int_fg, rune_bg)
				x_copy++
				if x_copy > width-2 {
					x_copy = 2
					y_copy += 3
				}
				if y_copy > height-2 {
					break
				}
			}
		}
		str := fmt.Sprintf("%02x", b)
		x += drawStringAtPoint(str, x, y, hex_fg, hex_bg)
		x++
	}
	termbox.Flush()

	return new_view_port
}

func drawColorScreen(fg termbox.Attribute, bg termbox.Attribute) {
	width, height := termbox.Size()
	drawBackground(bg)
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
	termbox.Flush()
}

func mainLoop(bytes []byte, style Style) {
	var view_port ViewPort
	var cursor Cursor
	cursor.int_length = 4
	cursor.fp_length = 4
	cursor.mode = StringMode

	hilite := cursor.highlightRange(bytes)
	view_port = drawDataScreen(bytes, view_port, style, cursor, hilite)
	prev_mode := cursor.mode
	display_screen := DataScreen
	modes := map[rune]CursorMode{
		'i': IntegerMode,
		't': StringMode,
		'f': FloatingPointMode,
		'p': BitPatternMode,
	}
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlP { /* color palette */
				if display_screen == ColorScreen {
					display_screen = DataScreen
				} else {
					display_screen = ColorScreen
				}
			} else if event.Ch == '?' { /* about */
				if display_screen == AboutScreen {
					display_screen = DataScreen
				} else {
					display_screen = AboutScreen
				}
			} else if display_screen == ColorScreen || display_screen == AboutScreen {
				display_screen = DataScreen
			} else if event.Ch == 'j' || event.Key == termbox.KeyArrowDown { /* down */
				cursor.pos += view_port.bytes_per_row
			} else if event.Key == termbox.KeyCtrlF || event.Key == termbox.KeyPgdn { /* page down */
				cursor.pos += view_port.bytes_per_row * view_port.number_of_rows
			} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp { /* up */
				cursor.pos -= view_port.bytes_per_row
			} else if event.Key == termbox.KeyCtrlB || event.Key == termbox.KeyPgup { /* page up */
				cursor.pos -= view_port.bytes_per_row * view_port.number_of_rows
			} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft { /* left */
				cursor.pos--
			} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight { /* right */
				cursor.pos++
			} else if event.Ch == 'w' { /* forward 1 "word" */
				cursor.pos += 4
			} else if event.Ch == 'b' { /* back 1 "word" */
				cursor.pos -= 4
			} else if modes[event.Ch] != 0 {
				if cursor.mode == modes[event.Ch] {
					cursor.mode = prev_mode
					prev_mode = modes[event.Ch]
				} else {
					prev_mode = cursor.mode
					cursor.mode = modes[event.Ch]
				}
			} else if event.Ch == 'u' || event.Ch == 'U' {
				if cursor.mode == IntegerMode {
					cursor.unsigned = !cursor.unsigned
				}
			} else if event.Ch == 'e' || event.Ch == 'E' {
				if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
					cursor.big_endian = !cursor.big_endian
				}
			} else if event.Ch == 'H' { /* shorten */
				if cursor.mode == IntegerMode && cursor.int_length > MIN_INTEGER_WIDTH {
					cursor.int_length /= 2
				}
				if cursor.mode == FloatingPointMode && cursor.fp_length > MIN_FLOATING_POINT_WIDTH {
					cursor.fp_length /= 2
				}
			} else if event.Ch == 'L' || event.Ch == ':' { /* lengthen */
				if cursor.mode == IntegerMode && cursor.int_length < MAX_INTEGER_WIDTH {
					cursor.int_length *= 2
				}
				if cursor.mode == FloatingPointMode && cursor.fp_length < MAX_FLOATING_POINT_WIDTH {
					cursor.fp_length *= 2
				}
			} else if event.Key == termbox.KeyCtrlE { /* scroll down */
				view_port.first_row++
				if cursor.pos < view_port.first_row*view_port.bytes_per_row {
					cursor.pos += view_port.bytes_per_row
				}
			} else if event.Key == termbox.KeyCtrlY { /* scroll up */
				view_port.first_row--
				if cursor.pos > (view_port.first_row+view_port.number_of_rows)*view_port.bytes_per_row {
					cursor.pos -= view_port.bytes_per_row
				}
			} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
				break
			}
			if cursor.pos < 0 {
				cursor.pos = 0
			}
			if cursor.pos+cursor.length() > len(bytes) {
				cursor.pos = len(bytes) - cursor.length()
			}
			if cursor.pos >= (view_port.first_row+view_port.number_of_rows)*view_port.bytes_per_row {
				view_port.first_row += view_port.number_of_rows
			}
			if cursor.pos < view_port.first_row*view_port.bytes_per_row {
				if view_port.first_row >= view_port.number_of_rows {
					view_port.first_row -= view_port.number_of_rows
				} else {
					view_port.first_row = 0
				}
			}
			if display_screen == DataScreen {
				hilite = cursor.highlightRange(bytes)
				view_port = drawDataScreen(bytes, view_port, style, cursor, hilite)
			} else if display_screen == ColorScreen {
				drawColorScreen(style.default_fg, style.default_bg)
			} else if display_screen == AboutScreen {
				drawAboutScreen(style.default_fg, style.default_bg)
			}
		}
		if event.Type == termbox.EventResize {
			if display_screen == DataScreen {
				view_port = drawDataScreen(bytes, view_port, style, cursor, hilite)
			} else if display_screen == ColorScreen {
				drawColorScreen(style.default_fg, style.default_bg)
			} else if display_screen == AboutScreen {
				drawAboutScreen(style.default_fg, style.default_bg)
			}
		}
		if event.Type == termbox.EventMouse {
		}
	}
}

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	path := os.Args[1]

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %q\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes from %s\n", len(bytes), path)
	if len(bytes) < 8 {
		fmt.Printf("File %s is too short to be edited\n", path)
		os.Exit(1)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var style Style
	style.default_bg = termbox.Attribute(1)
	style.default_fg = termbox.Attribute(256)
	style.rune_fg = termbox.Attribute(248)
	style.int_fg = termbox.Attribute(154)
	style.bit_fg = termbox.Attribute(154)
	style.space_rune_fg = termbox.Attribute(240)

	style.text_cursor_hex_bg = termbox.Attribute(167)
	style.bit_cursor_hex_bg = termbox.Attribute(26)
	style.int_cursor_hex_bg = termbox.Attribute(63)
	style.fp_cursor_hex_bg = termbox.Attribute(127)

	style.hilite_hex_fg = termbox.Attribute(231)
	style.hilite_rune_fg = termbox.Attribute(256)

	termbox.SetOutputMode(termbox.Output256)

	mainLoop(bytes, style)
}
