package main

import "fmt"

type OffsetWidget int

func (widget OffsetWidget) sizeForLayout(layout Layout) Size {
	if layout.pressure > 3 {
		return Size{0, 0}
	}
	height := 2
	width := 20
	if layout.show_date {
		height = 4
	}
	return Size{width, height}
}

func (widget OffsetWidget) drawAtPoint(tab *DataTab, layout Layout, point Point, style Style) Size {
	if layout.pressure > 3 {
		return Size{0, 0}
	}
	fg := style.default_fg
	bg := style.default_bg
	cursor := tab.cursor
	y_pos := point.y
	x_pos := point.x
	width := 20
	if tab.edit_mode == EditingSearch || tab.is_searching {
		x_pos += drawStringAtPoint("Search(/)", point.x, y_pos, fg, bg)
		if tab.is_searching {
			drawStringAtPoint(tab.prev_search, x_pos+2, y_pos, fg, bg)
		}
	} else if cursor.hex_mode {
		last_pos := len(tab.bytes) - 1
		last_pos_len := len(fmt.Sprintf("%x", last_pos))
		drawStringAtPoint(fmt.Sprintf("Offset(:)  0x%0[2]*[1]x", cursor.pos, last_pos_len), point.x, y_pos, fg, bg)
	} else {
		drawStringAtPoint(fmt.Sprintf("Offset(:)  %d", cursor.pos), point.x, y_pos, fg, bg)
	}
	y_pos++
	x_pos = point.x
	if tab.is_searching {
		x_pos += drawStringAtPoint("[", x_pos, y_pos, fg, bg)
		eighths := [...]string{
			" ",
			"▏",
			"▎",
			"▍",
			"▌",
			"▋",
			"▊",
			"▉",
			"█",
		}
		fifty_sixths := int(7 * 8 * tab.search_progress)
		if fifty_sixths < 2 {
			drawStringAtPoint(fmt.Sprintf("%2.2f%% ", 100*tab.search_progress), x_pos+1, y_pos, style.space_rune_fg, bg)
		}
		for i := 0; i < 7; i++ {
			if fifty_sixths >= 8*(i+1) {
				drawStringAtPoint(eighths[8], x_pos+i, y_pos, style.search_progress_fg, bg)
			} else if fifty_sixths > 8*i {
				drawStringAtPoint(eighths[fifty_sixths-8*i], x_pos+i, y_pos, style.search_progress_fg, bg)
			}
		}
		x_pos += 7
		x_pos += drawStringAtPoint("]", x_pos, y_pos, fg, bg)
		drawStringAtPoint("<Esc> cancel", x_pos+2, y_pos, fg, bg)
	} else {
		drawStringAtPoint(fmt.Sprintf("  Type :  %s", cursor.c_type()), point.x, y_pos, fg, bg)
	}

	if layout.show_date {
		y_pos++
		y_pos++
		if cursor.mode == FloatingPointMode || cursor.mode == IntegerMode {
			fg = style.default_fg
		} else {
			fg = style.space_rune_fg
		}
		epoch_string := fmt.Sprintf(" Epoch(@) %s", cursor.epoch_time.Format("1/2/2006"))
		drawStringAtPoint(epoch_string, point.x, y_pos, fg, bg)
	}

	return Size{width, y_pos - point.y + 1}
}
