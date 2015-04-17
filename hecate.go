package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"unsafe"
)

type CursorMode int

const (
	StringMode CursorMode = iota
	BitMode
	IntegerMode
	FloatingPointMode
)

const MAX_INTEGER_WIDTH = 8
const MIN_INTEGER_WIDTH = 1
const MAX_FLOATING_POINT_WIDTH = 8
const MIN_FLOATING_POINT_WIDTH = 4

type ByteRange struct {
	pos    int
	length int
}

type Cursor struct {
	pos        int
	int_length int
	fp_length  int
	mode       CursorMode
	unsigned   bool
	big_endian bool
}

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

func isASCII(val byte) bool {
	return (val >= 0x20 && val < 0x80)
}

func isCode(val byte) bool {
	return val == 0x09 || val == 0x0A || val == 0x0D
}

func isPrintable(val byte) bool {
	return isASCII(val) || isCode(val)
}

func highlightRange(data []byte, cursor Cursor) ByteRange {
	var hilite ByteRange
	if cursor.mode != StringMode || !isPrintable(data[cursor.pos]) {
		return hilite
	}
	hilite.pos = cursor.pos
	for ; hilite.pos > 0 && isPrintable(data[hilite.pos-1]); hilite.pos-- {
	}
	for ; hilite.pos+hilite.length < len(data) && isPrintable(data[hilite.pos+hilite.length]); hilite.length++ {
	}
	return hilite
}

func cursorLength(cursor Cursor) int {
	if cursor.mode == IntegerMode {
		return cursor.int_length
	}
	if cursor.mode == FloatingPointMode {
		return cursor.fp_length
	}
	return 1
}

func cursorColor(cursor Cursor, style Style) termbox.Attribute {
	if cursor.mode == IntegerMode {
		return style.int_cursor_hex_bg
	}
	if cursor.mode == FloatingPointMode {
		return style.fp_cursor_hex_bg
	}
	if cursor.mode == BitMode {
		return style.bit_cursor_hex_bg
	}
	return style.text_cursor_hex_bg
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

func drawBytes(data []byte, old_view_port ViewPort, style Style, cursor Cursor, hilite ByteRange) ViewPort {
	x, y := 2, 1
	width, height := termbox.Size()
	rows := 1
	var new_view_port ViewPort
	new_view_port.bytes_per_row = (width - 3) / 3
	new_view_port.number_of_rows = (height - 3) / 3

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

	drawBackground(style.default_bg)
	start := new_view_port.first_row * new_view_port.bytes_per_row
	end := start + new_view_port.number_of_rows*new_view_port.bytes_per_row
	for index := start; index < end && index < len(data); index++ {
		b := data[index]
		hex_fg := style.default_fg
		hex_bg := style.default_bg
		code_fg := style.space_rune_fg
		rune_fg := style.rune_fg
		rune_bg := style.default_bg
		cursor_length := cursorLength(cursor)
		if x+3 > width-1 {
			x = 2
			y += 3
			rows++
		}
		if y > height-2 {
			break
		}
		if index >= cursor.pos && index < cursor.pos+cursor_length {
			hex_bg = cursorColor(cursor, style)
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
		} else if cursor.mode == BitMode {
			for i := 0; i < 8; i++ {
				if b&(1<<uint8(7-i)) > 0 {
					termbox.SetCell(x-1+(i%4), y+1+i/4, '●', style.bit_fg, rune_bg)
				} else {
					termbox.SetCell(x-1+(i%4), y+1+i/4, '○', style.bit_fg, rune_bg)
				}
			}
		} else if index == cursor.pos {
			str := ""
			cursor_length := cursorLength(cursor)
			total_length := cursor_length*3 + 1
			var integer uint64
			if cursor.big_endian {
				for i := 0; i < cursor_length; i++ {
					integer = (integer * 256) + uint64(data[cursor.pos+i])
				}
			} else {
				for i := cursor_length - 1; i >= 0; i-- {
					integer = (integer * 256) + uint64(data[cursor.pos+i])
				}
			}
			if cursor.mode == IntegerMode {
				if cursor.int_length == 1 {
					if cursor.unsigned {
						str = fmt.Sprintf("%d", uint8(integer))
					} else {
						str = fmt.Sprintf("%d", int8(integer))
					}
				} else if cursor.int_length == 2 {
					if cursor.unsigned {
						str = fmt.Sprintf("%d", uint16(integer))
					} else {
						str = fmt.Sprintf("%d", int16(integer))
					}
				} else if cursor.int_length == 4 {
					if cursor.unsigned {
						str = fmt.Sprintf("%d", uint32(integer))
					} else {
						str = fmt.Sprintf("%d", int32(integer))
					}
				} else if cursor.int_length == 8 {
					if cursor.unsigned {
						str = fmt.Sprintf("%d", uint64(integer))
					} else {
						str = fmt.Sprintf("%d", int64(integer))
					}
				}
			} else if cursor.mode == FloatingPointMode {
				if cursor.fp_length == 4 {
					var integer32 uint32 = uint32(integer)
					str = fmt.Sprintf("%.5g", *(*float32)(unsafe.Pointer(&integer32)))
				} else if cursor.fp_length == 8 {
					str = fmt.Sprintf("%g", *(*float64)(unsafe.Pointer(&integer)))
				}
			}
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
			}
		}
		str := fmt.Sprintf("%02x", b)
		for _, runeValue := range str {
			termbox.SetCell(x, y, runeValue, hex_fg, hex_bg)
			x++
		}
		x++
	}
	termbox.Flush()

	return new_view_port
}

func drawColors(fg termbox.Attribute, bg termbox.Attribute) {
	width, height := termbox.Size()
	drawBackground(bg)
	x, y := 1, 1
	for color := 1; color <= 256; color++ {
		if x+8 > width {
			x = 1
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
		for _, runeValue := range str {
			termbox.SetCell(x, y, runeValue, fg, bg)
			x++
		}
		x += 2
	}
	termbox.Flush()
}

func main() {
	var err error
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

	var cursor Cursor
	cursor.int_length = 4
	cursor.fp_length = 4

	var hilite ByteRange

	var view_port ViewPort
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
	hilite = highlightRange(bytes, cursor)
	view_port = drawBytes(bytes, view_port, style, cursor, hilite)
	prev_mode := cursor.mode
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Ch == 'j' || event.Key == termbox.KeyArrowDown {
				cursor.pos += view_port.bytes_per_row
			} else if event.Key == termbox.KeyCtrlF || event.Key == termbox.KeyPgdn {
				cursor.pos += view_port.bytes_per_row * view_port.number_of_rows
			} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp {
				cursor.pos -= view_port.bytes_per_row
			} else if event.Key == termbox.KeyCtrlB || event.Key == termbox.KeyPgup {
				cursor.pos -= view_port.bytes_per_row * view_port.number_of_rows
			} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft {
				if cursor.pos > 0 {
					cursor.pos--
				}
			} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight {
				if cursor.pos < len(bytes)-1 {
					cursor.pos++
				}
			} else if event.Ch == 'w' {
				if cursor.pos+4 < len(bytes) {
					cursor.pos += 4
				} else {
					cursor.pos = len(bytes) - 1
				}
			} else if event.Ch == 'b' {
				if cursor.pos-4 >= 0 {
					cursor.pos -= 4
				} else {
					cursor.pos = 0
				}
			} else if event.Ch == 'i' {
				if cursor.mode == IntegerMode {
					cursor.mode = prev_mode
					prev_mode = IntegerMode
				} else {
					prev_mode = cursor.mode
					cursor.mode = IntegerMode
				}
			} else if event.Ch == 'u' {
				if cursor.mode == IntegerMode {
					cursor.unsigned = !cursor.unsigned
				}
			} else if event.Ch == 'e' {
				if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
					cursor.big_endian = !cursor.big_endian
				}
			} else if event.Ch == 'f' {
				if cursor.mode == FloatingPointMode {
					cursor.mode = prev_mode
					prev_mode = FloatingPointMode
				} else {
					prev_mode = cursor.mode
					cursor.mode = FloatingPointMode
				}
			} else if event.Ch == 'm' {
				if cursor.mode == BitMode {
					cursor.mode = prev_mode
					prev_mode = BitMode
				} else {
					prev_mode = cursor.mode
					cursor.mode = BitMode
				}
			} else if event.Ch == 't' {
				if cursor.mode == StringMode {
					cursor.mode = prev_mode
					prev_mode = StringMode
				} else {
					prev_mode = cursor.mode
					cursor.mode = StringMode
				}
			} else if event.Key == termbox.KeyCtrlH { /* shorten */
				if cursor.mode == IntegerMode && cursor.int_length > MIN_INTEGER_WIDTH {
					cursor.int_length /= 2
				}
				if cursor.mode == FloatingPointMode && cursor.fp_length > MIN_FLOATING_POINT_WIDTH {
					cursor.fp_length /= 2
				}
			} else if event.Key == termbox.KeyCtrlL { /* lengthen */
				if cursor.mode == IntegerMode && cursor.int_length < MAX_INTEGER_WIDTH {
					cursor.int_length *= 2
				}
				if cursor.mode == FloatingPointMode && cursor.fp_length < MAX_FLOATING_POINT_WIDTH {
					cursor.fp_length *= 2
				}
			} else if event.Key == termbox.KeyCtrlC {
				drawColors(style.default_fg, style.default_bg)
			} else {
				break
			}
			if cursor.pos < 0 {
				cursor.pos = 0
			}
			if cursor.pos+cursorLength(cursor) > len(bytes) {
				cursor.pos = len(bytes) - cursorLength(cursor)
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
			if event.Key != termbox.KeyCtrlC {
				hilite = highlightRange(bytes, cursor)
				view_port = drawBytes(bytes, view_port, style, cursor, hilite)
			}
		}
		if event.Type == termbox.EventResize {
			view_port = drawBytes(bytes, view_port, style, cursor, hilite)
		}
		if event.Type == termbox.EventMouse {
		}
	}
}
