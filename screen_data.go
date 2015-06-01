package main

import (
	"fmt"
	"math"
	"time"

	"github.com/nsf/termbox-go"
)

type EditMode int

const (
	EditingOffset EditMode = iota + 1
	EditingSearch
	EditingEpoch
)

type ViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type DataScreen struct {
	bytes                   []byte
	cursor                  Cursor
	hilite                  ByteRange
	view_port               ViewPort
	prev_mode               CursorMode
	prev_search             string
	edit_mode               EditMode
	show_date               bool
	is_searching            bool
	search_progress         float64
	search_progress_channel chan int
	search_result_channel   chan *Cursor
	search_quit_channel     chan bool
	field_editor            *FieldEditor
}

func scanEpoch(value string, epoch time.Time) time.Time {
	parsed_time, parse_error := time.Parse("1/2/2006", value)
	if parse_error != nil {
		return epoch
	}
	return parsed_time
}

func scanOffset(value string, file_pos int) int {
	var err error
	expression := ""
	scanned_file_pos := 0
	if n, _ := fmt.Sscanf(value, "+%s", &expression); n > 0 {
		if scanned_file_pos, err = evaluateExpression(expression); err == nil {
			return file_pos + scanned_file_pos
		}
		return -1
	}
	if scanned_file_pos, err = evaluateExpression(value); err == nil {
		if scanned_file_pos < 0 {
			if scanned_file_pos+file_pos < 0 {
				return 0
			}
			return scanned_file_pos + file_pos
		}
		return scanned_file_pos
	}
	return -1
}

func scanSearchString(value string, bytes []byte, cursor Cursor, quit <-chan bool, progress chan<- int) *Cursor {
	representations := make(map[string]*Cursor)

	var scanned_fp float64
	var rest_of_value string
	if n, _ := fmt.Sscanf(value, "%g%s", &scanned_fp, &rest_of_value); n > 0 && len(rest_of_value) == 0 {
		fp32_string := binaryStringForInteger32(math.Float32bits(float32(scanned_fp)), cursor.big_endian)
		fp32_cursor := Cursor{mode: FloatingPointMode, fp_length: 4, big_endian: cursor.big_endian}
		representations[fp32_string] = &fp32_cursor

		fp64_string := binaryStringForInteger64(math.Float64bits(scanned_fp), cursor.big_endian)
		fp64_cursor := Cursor{mode: FloatingPointMode, fp_length: 8, big_endian: cursor.big_endian}
		representations[fp64_string] = &fp64_cursor

		var scanned_int int64
		if n, _ := fmt.Sscanf(value, "%d%s", &scanned_int, &rest_of_value); n > 0 && scanned_fp == float64(scanned_int) && len(rest_of_value) == 0 {
			if scanned_int >= math.MinInt8 && scanned_int <= math.MaxUint8 {
				int8_string := binaryStringForInteger8(uint8(scanned_int))
				int8_cursor := Cursor{mode: IntegerMode, int_length: 1, unsigned: (scanned_int > math.MaxInt8)}
				representations[int8_string] = &int8_cursor
			}
			if scanned_int >= math.MinInt16 && scanned_int <= math.MaxUint16 {
				int16_string := binaryStringForInteger16(uint16(scanned_int), cursor.big_endian)
				int16_cursor := Cursor{mode: IntegerMode, int_length: 2, unsigned: (scanned_int > math.MaxInt16),
					big_endian: cursor.big_endian}
				representations[int16_string] = &int16_cursor
			}
			if scanned_int >= math.MinInt32 && scanned_int <= math.MaxUint32 {
				int32_string := binaryStringForInteger32(uint32(scanned_int), cursor.big_endian)
				int32_cursor := Cursor{mode: IntegerMode, int_length: 4, unsigned: (scanned_int > math.MaxInt32),
					big_endian: cursor.big_endian}
				representations[int32_string] = &int32_cursor
			}
			int64_string := binaryStringForInteger64(uint64(scanned_int), cursor.big_endian)
			int64_cursor := Cursor{mode: IntegerMode, int_length: 8, unsigned: (scanned_int > math.MaxInt64),
				big_endian: cursor.big_endian}
			representations[int64_string] = &int64_cursor
		}
	}
	text_cursor := Cursor{mode: StringMode}
	representations[value] = &text_cursor

	first_match := -1
	first_length := 1
	first_cursor := cursor

	for k, v := range representations {
		start_pos := cursor.pos + 1
		found_pos := -1
		if start_pos < len(bytes) {
			found_pos = interruptibleSearch(bytes[start_pos:], k, quit, progress)
			if found_pos >= 0 {
				found_pos += start_pos
			}
		}
		if found_pos == -1 {
			found_pos = interruptibleSearch(bytes[0:cursor.pos], k, quit, progress)
		}
		if found_pos == -2 {
			return nil
		}
		v.pos = found_pos
	}

	found_match := false

	for _, v := range representations {
		if v.pos != -1 {
			found_match = true
			ranked_pos := v.pos - cursor.pos
			if ranked_pos < 0 {
				ranked_pos += len(bytes)
			}
			if first_match == -1 || ranked_pos < first_match ||
				(ranked_pos == first_match && v.length() > first_length) {
				first_match = ranked_pos
				first_length = v.length()
				first_cursor.pos = v.pos
				if v.mode == FloatingPointMode {
					first_cursor.fp_length = v.fp_length
				} else if v.mode == IntegerMode {
					first_cursor.int_length = v.int_length
				}
				first_cursor.mode = v.mode
				first_cursor.unsigned = v.unsigned
			}
		}
	}
	if !found_match {
		return nil
	}
	return &first_cursor
}

func (screen *DataScreen) initializeWithBytes(bytes []byte) {
	cursor := Cursor{int_length: 4, fp_length: 4, mode: StringMode,
		epoch_unit: SecondsSinceEpoch, epoch_time: time.Unix(0, 0).UTC()}

	screen.search_result_channel = make(chan *Cursor)
	screen.search_quit_channel = make(chan bool)
	screen.search_progress_channel = make(chan int)
	screen.bytes = bytes
	screen.cursor = cursor
	screen.hilite = cursor.highlightRange(bytes)
	screen.prev_mode = cursor.mode
}

func (screen *DataScreen) receiveEvents(input <-chan termbox.Event, output chan<- int, quit <-chan bool) {
	for {
		do_quit := false
		select {
		case event := <-input:
			output <- screen.handleKeyEvent(event)
		case bytes_read := <-screen.search_progress_channel:
			screen.search_progress += float64(bytes_read) / float64(len(screen.bytes))
			if screen.search_progress > 1.0 {
				screen.search_progress = 0.0
			}
			output <- DATA_SCREEN_INDEX
		case search_result := <-screen.search_result_channel:
			screen.is_searching = false
			if search_result != nil {
				screen.cursor = *search_result
				screen.hilite = screen.cursor.highlightRange(screen.bytes)
			}
			output <- DATA_SCREEN_INDEX
		case <-quit:
			do_quit = true
		}
		if do_quit {
			if screen.is_searching {
				screen.search_quit_channel <- true
			}
			break
		}
	}
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
	} else if screen.field_editor != nil {
		new_pos := -1
		string_value, is_done := screen.field_editor.handleKeyEvent(event)
		if is_done {
			if len(string_value) > 0 {
				if screen.edit_mode == EditingSearch {
					screen.is_searching = true
					screen.search_progress = 0.0
					screen.prev_search = string_value
					go func() {
						cursor := scanSearchString(string_value, screen.bytes, screen.cursor,
							screen.search_quit_channel, screen.search_progress_channel)
						screen.search_result_channel <- cursor
					}()
				} else if screen.edit_mode == EditingOffset {
					new_pos = scanOffset(string_value, screen.cursor.pos)
				} else if screen.edit_mode == EditingEpoch {
					screen.cursor.epoch_time = scanEpoch(string_value, screen.cursor.epoch_time)
				}
			}
			screen.edit_mode = 0
			screen.field_editor = nil
		}
		if new_pos >= 0 {
			screen.cursor.pos = new_pos
		}
	} else if event.Ch == 'n' {
		if len(screen.prev_search) > 0 && !screen.is_searching {
			go func() {
				cursor := scanSearchString(screen.prev_search, screen.bytes, screen.cursor,
					screen.search_quit_channel, screen.search_progress_channel)
				screen.search_result_channel <- cursor
			}()
			screen.is_searching = true
		}
	} else if event.Ch == ':' {
		if screen.is_searching {
			screen.search_quit_channel <- true
		}
		screen.field_editor = new(FieldEditor)
		screen.edit_mode = EditingOffset
	} else if event.Ch == '/' {
		if screen.is_searching {
			screen.search_quit_channel <- true
		}
		screen.field_editor = new(FieldEditor)
		screen.edit_mode = EditingSearch
	} else if event.Ch == '@' {
		if screen.show_date {
			screen.field_editor = new(FieldEditor)
			screen.edit_mode = EditingEpoch
		}
	} else if event.Ch == 'x' {
		screen.cursor.hex_mode = !screen.cursor.hex_mode
	} else if event.Ch == 'D' {
		screen.show_date = !screen.show_date
	} else if event.Ch == 's' {
		if screen.show_date {
			screen.cursor.epoch_unit = SecondsSinceEpoch
		}
	} else if event.Ch == 'd' {
		if screen.show_date {
			screen.cursor.epoch_unit = DaysSinceEpoch
		}
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
	} else if event.Ch == '^' {
		screen.cursor.pos = screen.cursor.pos / screen.view_port.bytes_per_row * screen.view_port.bytes_per_row
	} else if event.Ch == '$' {
		screen.cursor.pos = (screen.cursor.pos/screen.view_port.bytes_per_row+1)*screen.view_port.bytes_per_row - screen.cursor.length()
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
	} else if event.Key == termbox.KeyCtrlD {
		if screen.is_searching {
			screen.search_quit_channel <- true
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
	screen.hilite = screen.cursor.highlightRange(screen.bytes)
	if screen.field_editor == nil {
		termbox.HideCursor()
	}

	return DATA_SCREEN_INDEX
}

func (screen *DataScreen) performLayout() {
	cursor := screen.cursor
	width, height := termbox.Size()
	legend_height := heightOfWidgets(screen.show_date)
	line_height := 3
	cursor_row_within_view_port := 0

	if cursor.pos >= (screen.view_port.first_row+screen.view_port.number_of_rows)*screen.view_port.bytes_per_row {
		screen.view_port.first_row += screen.view_port.number_of_rows
	}
	for cursor.pos < screen.view_port.first_row*screen.view_port.bytes_per_row {
		screen.view_port.first_row -= screen.view_port.number_of_rows
	}

	var new_view_port ViewPort
	new_view_port.bytes_per_row = (width - 3) / 3
	new_view_port.number_of_rows = (height - 1 - legend_height) / line_height
	new_view_port.first_row = screen.view_port.first_row

	if screen.view_port.bytes_per_row > 0 {
		cursor_row_within_view_port = cursor.pos/screen.view_port.bytes_per_row - screen.view_port.first_row
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

	screen.view_port = new_view_port
}

func (screen *DataScreen) drawScreen(style Style) {
	cursor := screen.cursor
	hilite := screen.hilite
	view_port := screen.view_port

	layout := drawWidgets(screen, style)
	x, y := 2, 1
	x_pad := 2
	line_height := 3
	width, height := termbox.Size()

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

	if screen.field_editor != nil {
		widget_width := layout.width()
		widget_height := layout.widget_size.height
		if layout.pressure < 4 {
			x = (width-widget_width)/2 + widget_width - 11
			if screen.edit_mode == EditingEpoch {
				y = height - 1
			} else {
				y = height - widget_height
			}
		} else {
			x = (width - 10) / 2
			y = height - widget_height - 1
		}
		termbox.SetCursor(x+2+screen.field_editor.cursor_pos, y)
		drawStringAtPoint(fmt.Sprintf(" %-10s ", screen.field_editor.value), x+1, y,
			style.field_editor_fg, style.field_editor_bg)
	}
}
