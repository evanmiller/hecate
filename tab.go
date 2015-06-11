package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

var modes = map[rune]CursorMode{
	'i': IntegerMode,
	't': StringMode,
	'f': FloatingPointMode,
	'p': BitPatternMode,
}

const (
	EditingOffset EditMode = iota + 1
	EditingSearch
	EditingEpoch
)

type EditMode int

type DataViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type DataTab struct {
	filename                string
	bytes                   []byte
	cursor                  Cursor
	hilite                  ByteRange
	view_port               DataViewPort
	prev_mode               CursorMode
	prev_search             string
	edit_mode               EditMode
	show_date               bool
	is_searching            bool
	search_progress         float64
	search_progress_channel chan int
	search_result_channel   chan *Cursor
	search_quit_channel     chan bool
	quit_channel            chan bool
	field_editor            *FieldEditor
}

func NewDataTab(file FileInfo) DataTab {
	cursor := Cursor{int_length: 4, fp_length: 4, bit_length: 1, mode: StringMode,
		epoch_unit: SecondsSinceEpoch, epoch_time: time.Unix(0, 0).UTC()}

	return DataTab{
		search_result_channel:   make(chan *Cursor),
		search_quit_channel:     make(chan bool),
		search_progress_channel: make(chan int),
		quit_channel:            make(chan bool, 10),
		bytes:                   file.bytes,
		filename:                file.filename,
		cursor:                  cursor,
		hilite:                  cursor.highlightRange(file.bytes),
		prev_mode:               cursor.mode,
	}
}

func (tab *DataTab) receiveEvents(output chan<- int) {
	for {
		do_quit := false
		select {
		case bytes_read := <-tab.search_progress_channel:
			tab.search_progress += float64(bytes_read) / float64(len(tab.bytes))
			if tab.search_progress > 1.0 {
				tab.search_progress = 0.0
			}
			output <- DATA_SCREEN_INDEX
		case search_result := <-tab.search_result_channel:
			tab.is_searching = false
			tab.search_progress = 0.0
			if search_result != nil {
				tab.cursor = *search_result
				tab.hilite = tab.cursor.highlightRange(tab.bytes)
			}
			output <- DATA_SCREEN_INDEX
		case <-tab.quit_channel:
			do_quit = true
		}
		if do_quit {
			if tab.is_searching {
				tab.search_quit_channel <- true
			}
			break
		}
	}
}

func (tab *DataTab) performLayout(width int, height int) {
	cursor := tab.cursor
	legend_height := heightOfWidgets(tab.show_date)
	line_height := 3
	cursor_row_within_view_port := 0

	if cursor.pos >= (tab.view_port.first_row+tab.view_port.number_of_rows)*tab.view_port.bytes_per_row {
		tab.view_port.first_row += tab.view_port.number_of_rows
	}
	for cursor.pos < tab.view_port.first_row*tab.view_port.bytes_per_row {
		tab.view_port.first_row -= tab.view_port.number_of_rows
	}

	var new_view_port DataViewPort
	new_view_port.bytes_per_row = (width - 3) / 3
	new_view_port.number_of_rows = (height - 1 - legend_height) / line_height
	new_view_port.first_row = tab.view_port.first_row

	if tab.view_port.bytes_per_row > 0 {
		cursor_row_within_view_port = cursor.pos/tab.view_port.bytes_per_row - tab.view_port.first_row
		if cursor.pos/new_view_port.bytes_per_row > cursor_row_within_view_port {
			new_view_port.first_row = cursor.pos/new_view_port.bytes_per_row - cursor_row_within_view_port
		}
		if cursor.pos/new_view_port.bytes_per_row >= new_view_port.first_row+new_view_port.number_of_rows {
			new_view_port.first_row = cursor.pos/new_view_port.bytes_per_row - new_view_port.number_of_rows + 1
		}
	}
	if new_view_port.first_row < 0 {
		new_view_port.first_row = 0
	}

	tab.view_port = new_view_port
}

func (tab *DataTab) handleKeyEvent(event termbox.Event) int {
	if tab.field_editor != nil {
		new_pos := -1
		string_value, is_done := tab.field_editor.handleKeyEvent(event)
		if is_done {
			if len(string_value) > 0 {
				if tab.edit_mode == EditingSearch {
					tab.is_searching = true
					tab.search_progress = 0.0
					tab.prev_search = string_value
					go func() {
						cursor := scanSearchString(string_value, tab.bytes, tab.cursor,
							tab.search_quit_channel, tab.search_progress_channel)
						tab.search_result_channel <- cursor
					}()
				} else if tab.edit_mode == EditingOffset {
					new_pos = scanOffset(string_value, tab.cursor.pos)
				} else if tab.edit_mode == EditingEpoch {
					tab.cursor.epoch_time = scanEpoch(string_value, tab.cursor.epoch_time)
				}
			}
			tab.edit_mode = 0
			tab.field_editor = nil
		}
		if new_pos >= 0 {
			tab.cursor.pos = new_pos
		}
	} else if event.Ch == 'q' || event.Key == termbox.KeyCtrlC {
		if tab.is_searching {
			tab.search_quit_channel <- true
		} else {
			return EXIT_SCREEN_INDEX
		}
	} else if event.Key == termbox.KeyEsc {
		if tab.is_searching {
			tab.search_quit_channel <- true
		}
	} else if event.Ch == 'n' {
		if len(tab.prev_search) > 0 && !tab.is_searching {
			go func() {
				cursor := scanSearchString(tab.prev_search, tab.bytes, tab.cursor,
					tab.search_quit_channel, tab.search_progress_channel)
				tab.search_result_channel <- cursor
			}()
			tab.is_searching = true
		}
	} else if event.Ch == ':' {
		if tab.is_searching {
			tab.search_quit_channel <- true
		}
		tab.field_editor = new(FieldEditor)
		tab.edit_mode = EditingOffset
	} else if event.Ch == '/' {
		if tab.is_searching {
			tab.search_quit_channel <- true
		}
		tab.field_editor = new(FieldEditor)
		tab.edit_mode = EditingSearch
	} else if event.Ch == '@' {
		if tab.show_date {
			tab.field_editor = new(FieldEditor)
			tab.edit_mode = EditingEpoch
		}
	} else if event.Ch == 'x' {
		tab.cursor.hex_mode = !tab.cursor.hex_mode
	} else if event.Ch == 'a' {
		tab.show_date = !tab.show_date
	} else if event.Ch == 's' {
		if tab.show_date {
			tab.cursor.epoch_unit = SecondsSinceEpoch
		}
	} else if event.Ch == 'd' {
		if tab.show_date {
			tab.cursor.epoch_unit = DaysSinceEpoch
		}
	} else if event.Ch == 'j' || event.Key == termbox.KeyArrowDown { // down
		tab.cursor.pos += tab.view_port.bytes_per_row
	} else if event.Key == termbox.KeyCtrlF || event.Key == termbox.KeyPgdn { // page down
		tab.cursor.pos += tab.view_port.bytes_per_row * tab.view_port.number_of_rows
	} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp { // up
		tab.cursor.pos -= tab.view_port.bytes_per_row
	} else if event.Key == termbox.KeyCtrlB || event.Key == termbox.KeyPgup { // page up
		tab.cursor.pos -= tab.view_port.bytes_per_row * tab.view_port.number_of_rows
	} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft { // left
		tab.cursor.pos--
	} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight { // right
		tab.cursor.pos++
	} else if event.Ch == 'w' { /* forward 1 "word" */
		tab.cursor.pos += 4
	} else if event.Ch == 'b' { /* back 1 "word" */
		tab.cursor.pos -= 4
	} else if event.Ch == 'g' {
		tab.cursor.pos = 0
	} else if event.Ch == 'G' {
		tab.cursor.pos = len(tab.bytes)
	} else if event.Ch == '^' {
		tab.cursor.pos = tab.cursor.pos / tab.view_port.bytes_per_row * tab.view_port.bytes_per_row
	} else if event.Ch == '$' {
		tab.cursor.pos = (tab.cursor.pos/tab.view_port.bytes_per_row+1)*tab.view_port.bytes_per_row - tab.cursor.length()
	} else if modes[event.Ch] != 0 {
		if tab.cursor.mode == modes[event.Ch] {
			tab.cursor.mode = tab.prev_mode
			tab.prev_mode = modes[event.Ch]
		} else {
			tab.prev_mode = tab.cursor.mode
			tab.cursor.mode = modes[event.Ch]
		}
	} else if event.Ch == 'u' || event.Ch == 'U' {
		if tab.cursor.mode == IntegerMode {
			tab.cursor.unsigned = !tab.cursor.unsigned
		}
	} else if event.Ch == 'e' || event.Ch == 'E' {
		if tab.cursor.mode == IntegerMode || tab.cursor.mode == FloatingPointMode {
			tab.cursor.big_endian = !tab.cursor.big_endian
		}
	} else if event.Ch == 'H' {
		tab.cursor.shrink()
	} else if event.Ch == 'L' {
		tab.cursor.grow()
	} else if event.Key == termbox.KeyCtrlE { // scroll down
		if (tab.view_port.first_row+1)*tab.view_port.bytes_per_row < len(tab.bytes) {
			tab.view_port.first_row++
			if tab.cursor.pos < tab.view_port.first_row*tab.view_port.bytes_per_row {
				tab.cursor.pos += tab.view_port.bytes_per_row
			}
		}
	} else if event.Key == termbox.KeyCtrlY { /* scroll up */
		tab.view_port.first_row--
		if tab.cursor.pos > (tab.view_port.first_row+tab.view_port.number_of_rows)*tab.view_port.bytes_per_row {
			tab.cursor.pos -= tab.view_port.bytes_per_row
		}
	}
	if tab.cursor.pos < 0 {
		tab.cursor.pos = 0
	}
	if tab.cursor.pos+tab.cursor.length() > len(tab.bytes) {
		tab.cursor.pos = len(tab.bytes) - tab.cursor.length()
	}
	tab.hilite = tab.cursor.highlightRange(tab.bytes)
	if tab.field_editor == nil {
		termbox.HideCursor()
	}

	return DATA_SCREEN_INDEX
}

func (tab *DataTab) drawTab(style Style, vertical_offset int) {
	cursor := tab.cursor
	hilite := tab.hilite
	view_port := tab.view_port

	layout := drawWidgets(tab, style)
	start_x, start_y := 2, 1+vertical_offset
	x, y := start_x, start_y
	x_pad := 2
	line_height := 3
	width, height := termbox.Size()

	last_y := y + view_port.number_of_rows*line_height - 1
	last_x := x + view_port.bytes_per_row*3 - 1

	y -= line_height

	cursor_x := x
	cursor_y := y
	cursor_length := cursor.length()
	start := view_port.first_row * view_port.bytes_per_row
	end := start + view_port.number_of_rows*view_port.bytes_per_row
	rune_bg := style.default_bg
	for index := start; index < end && index < len(tab.bytes); index++ {
		b := tab.bytes[index]
		hex_fg := style.default_fg
		hex_bg := style.default_bg
		code_fg := style.space_rune_fg
		rune_fg := style.rune_fg
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
		} else if index == cursor.pos {
			cursor_x = x
			cursor_y = y
		}
		str := fmt.Sprintf("%02x", b)
		x += drawStringAtPoint(str, x, y, hex_fg, hex_bg)
		x++
	}

	if cursor.mode == BitPatternMode {
		if cursor_length == 1 || (cursor.pos+1)%view_port.bytes_per_row == 0 {
			for j := 0; j < cursor_length; j++ {
				b := tab.bytes[cursor.pos+j]
				for i := 0; i < 8; i++ {
					if b&(1<<uint8(7-i)) > 0 {
						termbox.SetCell(cursor_x-1+(i%4), cursor_y+1+i/4, '●', style.bit_fg, rune_bg)
					} else {
						termbox.SetCell(cursor_x-1+(i%4), cursor_y+1+i/4, '○', style.bit_fg, rune_bg)
					}
				}
				cursor_x = start_x
				cursor_y += line_height
				if cursor_y > last_y {
					break
				}
			}
		} else {
			for j := 0; j < cursor_length; j++ {
				b := tab.bytes[cursor.pos+j]
				for i := 0; i < 8; i++ {
					if b&(1<<uint8(7-i)) > 0 {
						termbox.SetCell(cursor_x-1+i, cursor_y+j+1, '●', style.bit_fg, rune_bg)
					} else {
						termbox.SetCell(cursor_x-1+i, cursor_y+j+1, '○', style.bit_fg, rune_bg)
					}
				}
			}
		}
	} else if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
		total_length := cursor_length*3 + 1
		str := cursor.formatBytesAsNumber(tab.bytes[cursor.pos : cursor.pos+cursor_length])
		x_copy := cursor_x - 1
		y_copy := cursor_y + 1
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

	if tab.field_editor != nil {
		widget_width := layout.width()
		widget_height := layout.widget_size.height
		if layout.pressure < 4 {
			x = (width-widget_width)/2 + widget_width - 11
			if tab.edit_mode == EditingEpoch {
				y = height - 1
			} else {
				y = height - widget_height
			}
		} else {
			x = (width - 10) / 2
			y = height - widget_height - 1
		}
		termbox.SetCursor(x+2+tab.field_editor.cursor_pos, y)
		drawStringAtPoint(fmt.Sprintf(" %-10s ", tab.field_editor.value), x+1, y,
			style.field_editor_fg, style.field_editor_bg)
	}
}
