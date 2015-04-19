package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type ViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type DataScreen struct {
	bytes     []byte
	cursor    Cursor
	hilite    ByteRange
	view_port ViewPort
	prev_mode CursorMode
}

func (screen *DataScreen) handleKeyEvent(event termbox.Event) int {
	modes := map[rune]CursorMode{
		'i': IntegerMode,
		't': StringMode,
		'f': FloatingPointMode,
		'p': BitPatternMode,
	}
	if event.Key == termbox.KeyCtrlP { // color palette
		return PALETTE_SCREEN_INDEX
	} else if event.Ch == '?' { // about
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'j' || event.Key == termbox.KeyArrowDown { // down
		screen.cursor.pos += screen.view_port.bytes_per_row
	} else if event.Key == termbox.KeyCtrlF || event.Key == termbox.KeyPgdn { // page down
		screen.cursor.pos += screen.view_port.bytes_per_row * screen.view_port.number_of_rows
	} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp { // up
		screen.cursor.pos -= screen.view_port.bytes_per_row
	} else if event.Key == termbox.KeyCtrlB || event.Key == termbox.KeyPgup { // page up
		screen.cursor.pos -= screen.view_port.bytes_per_row * screen.view_port.number_of_rows
	} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft { // left
		screen.cursor.pos--
	} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight { // right
		screen.cursor.pos++
	} else if event.Ch == 'w' { /* forward 1 "word" */
		screen.cursor.pos += 4
	} else if event.Ch == 'b' { /* back 1 "word" */
		screen.cursor.pos -= 4
	} else if event.Ch == 'g' {
		screen.cursor.pos = 0
	} else if event.Ch == 'G' {
		screen.cursor.pos = len(screen.bytes)
	} else if modes[event.Ch] != 0 {
		if screen.cursor.mode == modes[event.Ch] {
			screen.cursor.mode = screen.prev_mode
			screen.prev_mode = modes[event.Ch]
		} else {
			screen.prev_mode = screen.cursor.mode
			screen.cursor.mode = modes[event.Ch]
		}
	} else if event.Ch == 'u' || event.Ch == 'U' {
		if screen.cursor.mode == IntegerMode {
			screen.cursor.unsigned = !screen.cursor.unsigned
		}
	} else if event.Ch == 'e' || event.Ch == 'E' {
		if screen.cursor.mode == IntegerMode || screen.cursor.mode == FloatingPointMode {
			screen.cursor.big_endian = !screen.cursor.big_endian
		}
	} else if event.Ch == 'H' { /* shorten */
		if screen.cursor.length() > screen.cursor.minimumLength() {
			if screen.cursor.mode == IntegerMode {
				screen.cursor.int_length /= 2
			} else if screen.cursor.mode == FloatingPointMode {
				screen.cursor.fp_length /= 2
			}
		}
	} else if event.Ch == 'L' { /* lengthen */
		if screen.cursor.length() < screen.cursor.maximumLength() {
			if screen.cursor.mode == IntegerMode {
				screen.cursor.int_length *= 2
			} else if screen.cursor.mode == FloatingPointMode {
				screen.cursor.fp_length *= 2
			}
		}
	} else if event.Key == termbox.KeyCtrlE { // scroll down
		if (screen.view_port.first_row+1)*screen.view_port.bytes_per_row < len(screen.bytes) {
			screen.view_port.first_row++
			if screen.cursor.pos < screen.view_port.first_row*screen.view_port.bytes_per_row {
				screen.cursor.pos += screen.view_port.bytes_per_row
			}
		}
	} else if event.Key == termbox.KeyCtrlY { /* scroll up */
		screen.view_port.first_row--
		if screen.cursor.pos > (screen.view_port.first_row+screen.view_port.number_of_rows)*screen.view_port.bytes_per_row {
			screen.cursor.pos -= screen.view_port.bytes_per_row
		}
	} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
		return EXIT_SCREEN_INDEX
	}
	if screen.cursor.pos < 0 {
		screen.cursor.pos = 0
	}
	if screen.cursor.pos+screen.cursor.length() > len(screen.bytes) {
		screen.cursor.pos = len(screen.bytes) - screen.cursor.length()
	}
	if screen.cursor.pos >= (screen.view_port.first_row+screen.view_port.number_of_rows)*screen.view_port.bytes_per_row {
		screen.view_port.first_row += screen.view_port.number_of_rows
	}
	for screen.cursor.pos < screen.view_port.first_row*screen.view_port.bytes_per_row {
		screen.view_port.first_row -= screen.view_port.number_of_rows
		if screen.view_port.first_row < 0 {
			screen.view_port.first_row = 0
		}
	}
	screen.hilite = screen.cursor.highlightRange(screen.bytes)

	return DATA_SCREEN_INDEX
}

func (screen *DataScreen) performLayout() {
	width, height := termbox.Size()
	legend_height := heightOfWidgets()
	line_height := 3

	var new_view_port ViewPort
	new_view_port.bytes_per_row = (width - 3) / 3
	new_view_port.number_of_rows = (height - 1 - legend_height) / line_height

	cursor := screen.cursor
	cursor_row_within_view_port := 0
	if screen.view_port.bytes_per_row > 0 {
		cursor_row_within_view_port = cursor.pos/screen.view_port.bytes_per_row - screen.view_port.first_row
		if cursor.pos/new_view_port.bytes_per_row > cursor_row_within_view_port {
			new_view_port.first_row = cursor.pos/screen.view_port.bytes_per_row - cursor_row_within_view_port
		}
		if cursor.pos/new_view_port.bytes_per_row >= new_view_port.first_row+new_view_port.number_of_rows {
			new_view_port.first_row = cursor.pos/new_view_port.bytes_per_row - new_view_port.number_of_rows + 1
		}
	}

	screen.view_port = new_view_port
}

func (screen *DataScreen) drawScreen(style Style) {
	x, y := 2, 1
	x_pad := 2
	line_height := 3
	width, _ := termbox.Size()
	drawWidgets(screen.cursor, style)

	cursor := screen.cursor
	hilite := screen.hilite
	view_port := screen.view_port

	last_y := y + view_port.number_of_rows*line_height - 1
	last_x := x + view_port.bytes_per_row*3 - 1

	y = -2

	start := view_port.first_row * view_port.bytes_per_row
	end := start + view_port.number_of_rows*view_port.bytes_per_row
	for index := start; index < end && index < len(screen.bytes); index++ {
		b := screen.bytes[index]
		hex_fg := style.default_fg
		hex_bg := style.default_bg
		code_fg := style.space_rune_fg
		rune_fg := style.rune_fg
		rune_bg := style.default_bg
		cursor_length := cursor.length()
		if index%view_port.bytes_per_row == 0 {
			x = x_pad
			y += line_height
		}
		if y > last_y {
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
			total_length := cursor_length*3 + 1
			str := cursor.formatBytesAsNumber(screen.bytes[cursor.pos : cursor.pos+cursor_length])
			x_copy := x - 1
			y_copy := y + 1
			x_copy = x_copy + (total_length-len(str))/2
			if x_copy > last_x {
				x_copy = (x_copy % (width - x_pad)) + x_pad
				y_copy += line_height
			}
			for _, runeValue := range str {
				if y_copy > last_y {
					break
				}
				termbox.SetCell(x_copy, y_copy, runeValue, style.int_fg, rune_bg)
				x_copy++
				if x_copy > last_x {
					x_copy = x_pad
					y_copy += line_height
				}
			}
		}
		str := fmt.Sprintf("%02x", b)
		x += drawStringAtPoint(str, x, y, hex_fg, hex_bg)
		x++
	}
}
