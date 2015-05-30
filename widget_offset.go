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

func (widget OffsetWidget) drawAtPoint(cursor Cursor, data []byte, point Point, layout Layout, style Style, mode EditMode) Size {
	if layout.pressure > 3 {
		return Size{0, 0}
	}
	fg := style.default_fg
	bg := style.default_bg
	y_pos := point.y
	width := 20
	if mode == EditingSearch {
		drawStringAtPoint("Search(/)", point.x, y_pos, fg, bg)
	} else if cursor.hex_mode {
		drawStringAtPoint(fmt.Sprintf("Offset(:)  0x%x", cursor.pos), point.x, y_pos, fg, bg)
	} else {
		drawStringAtPoint(fmt.Sprintf("Offset(:)  %d", cursor.pos), point.x, y_pos, fg, bg)
	}
	y_pos++
	drawStringAtPoint(fmt.Sprintf("  Type :  %s", cursor.c_type()), point.x, y_pos, fg, bg)

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
